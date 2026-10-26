package agcapi

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/Createitv/agc-cli/pkg/domain"
)

func validateManagementContent(e domain.Endpoint, query map[string]string, body []byte, root map[string]any) []domain.InputViolation {
	if !IsManagementEndpoint(e) {
		return nil
	}
	var issues []domain.InputViolation
	add := func(path, message string) {
		issues = append(issues, domain.InputViolation{Code: "management_input", Field: path, Message: message, SourceURL: e.SourceURL, Suggestion: "Correct the field using " + e.Command + " --docs --pretty before invoking."})
	}
	text := func(key string) string { v, _ := root[key].(string); return v }
	if e.FamilyID == "cache" {
		if e.ID == "execute-command" {
			c := text("command")
			tokens := strings.Split(c, " ")
			supported := " DEL EXISTS EXPIRE PEXPIRE TTL TYPE APPEND DECR DECRBY GET GETSET INCR INCRBY MGET MSET MSETNX SET SETEX SETNX LINDEX LLEN LPOP LPUSH LRANGE LREM LSET LTRIM RPOP RPUSH SADD SCARD SISMEMBER SMEMBERS SPOP SRANDMEMBER SREM HDEL HEXISTS HGET HGETALL HINCRBY HLEN HKEYS HMGET HMSET HSET HVALS "
			valid := len(tokens) > 1 && len(c) <= 4096 && strings.Contains(supported, " "+strings.ToUpper(tokens[0])+" ")
			for _, b := range []byte(c) {
				if b < 32 || b > 126 {
					valid = false
				}
			}
			for _, token := range tokens {
				if token == "" {
					valid = false
				}
			}
			if !valid {
				add("body.command", "Use a supported Web CLI command with ASCII arguments, no empty or space-containing arguments, at most 4096 characters.")
			}
		}
		if e.ID == "reset-password" || e.ID == "set-password" {
			p := text("password")
			valid := len(p) >= 8 && len(p) <= 32 && regexp.MustCompile(`^[A-Za-z][A-Za-z0-9\-~!@#%^*_=+?$&()|<>{}\[\]]*$`).MatchString(p)
			for _, pattern := range []string{`[a-z]`, `[A-Z]`, `[0-9]`, `[\-~!@#%^*_=+?$&()|<>{}\[\]]`} {
				valid = valid && regexp.MustCompile(pattern).MatchString(p)
			}
			if !valid {
				add("body.password", "Password must start with a letter and contain lower/uppercase letters, a digit and an allowed special character, 8-32 ASCII characters.")
			}
		}
		if e.ID == "usage" {
			start, a := strconv.ParseInt(query["startTime"], 10, 64)
			end, b := strconv.ParseInt(query["endTime"], 10, 64)
			if a != nil || b != nil || start >= end || end-start > 259200000 {
				add("query.endTime", "Supply millisecond epoch startTime and endTime with startTime before endTime and an interval no longer than three days.")
			}
		}
	}
	if e.FamilyID == "functions" {
		if strings.Contains(e.ID, "alias") && (e.Method == "POST" || e.Method == "PUT") {
			config, ok := root["routingConfig"].(map[string]any)
			sum := 0
			valid := ok && len(config) > 0 && len(config) <= 2
			for version, value := range config {
				if version != "$latest" && !regexp.MustCompile(`^[1-9][0-9]{0,3}$`).MatchString(version) {
					valid = false
				}
				n, ok := value.(json.Number)
				weight, err := strconv.Atoi(n.String())
				if !ok || err != nil || weight < 0 || weight > 100 {
					valid = false
				}
				sum += weight
			}
			if !valid || sum != 100 {
				add("body.routingConfig", "Provide one or two version-to-integer-weight entries totaling 100; dynamic version keys are not literal required properties.")
			}
		}
		if strings.Contains(e.ID, "cron-trigger") {
			spec, _ := root["spec"].(map[string]any)
			schedule, _ := spec["schedule"].(string)
			valid := regexp.MustCompile(`^@every [1-9][0-9]*[smh]$`).MatchString(schedule)
			if strings.HasPrefix(schedule, "cron(") && strings.HasSuffix(schedule, ")") {
				valid = validCron(strings.TrimSuffix(strings.TrimPrefix(schedule, "cron("), ")"))
			}
			if !valid {
				add("body.spec.schedule", "Use @every <positive integer><s|m|h> or a five-field cron(...) expression with documented ranges.")
			}
		}
		if e.ID == "create-http-trigger" {
			if flag, ok := root["authFlag"].(bool); ok && !flag {
				add("body.authFlag", "Reviewed HTTP trigger creation requires authentication enabled (authFlag=true).")
			}
		}
		if e.ID == "create-layer-version" {
			encoded := text("zipFile")
			data, err := base64.StdEncoding.Strict().DecodeString(encoded)
			if err != nil || len(data) > 150<<20 {
				add("body.zipFile", "Supply Base64 ZIP bytes no larger than 150 MiB compressed.")
			} else {
				z, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
				if err != nil || len(z.File) > 30000 {
					add("body.zipFile", "ZIP must be valid and contain at most 30000 entries.")
				} else {
					var total uint64
					for _, f := range z.File {
						total += f.UncompressedSize64
						if total > 300<<20 {
							break
						}
						r, err := f.Open()
						if err != nil {
							total = 301 << 20
							break
						}
						n, err := io.Copy(io.Discard, io.LimitReader(r, int64(f.UncompressedSize64)+1))
						r.Close()
						if err != nil || uint64(n) != f.UncompressedSize64 {
							total = 301 << 20
							break
						}
					}
					if total > 300<<20 {
						add("body.zipFile", "ZIP content is corrupt or exceeds 300 MiB uncompressed.")
					}
				}
			}
		}
	}
	if e.FamilyID == "clouddb" {
		if (e.ID == "create-zone" || e.ID == "delete-zone") && strings.EqualFold(text("cloudDBZoneName"), "privatedefault") {
			add("body.cloudDBZoneName", "privatedefault is a reserved zone name.")
		}
		if e.ID == "import-object-types" && len(body) > 20<<20 {
			add("body.objectTypes", "Each schema import must be at most 20 MiB.")
		}
		for _, key := range []string{"objectType", "originObjectType"} {
			if model, ok := root[key].(map[string]any); ok {
				validateSchemaModel(model, "body."+key, add)
			}
		}
		if list, ok := root["objectTypes"].([]any); ok {
			for i, v := range list {
				if m, ok := v.(map[string]any); ok {
					validateSchemaModel(m, fmt.Sprintf("body.objectTypes[%d]", i), add)
				}
			}
		}
		if e.ID == "update-object-type" {
			before, _ := root["originObjectType"].(map[string]any)
			after, _ := root["objectType"].(map[string]any)
			if before["objectTypeName"] != after["objectTypeName"] {
				add("body.objectType.objectTypeName", "Object type names are immutable; use the exact prior name.")
			}
			oldFields := map[string]map[string]any{}
			oldList, _ := before["fields"].([]any)
			for _, v := range oldList {
				m, _ := v.(map[string]any)
				name, _ := m["fieldName"].(string)
				oldFields[strings.ToLower(name)] = m
			}
			newList, _ := after["fields"].([]any)
			for i, v := range newList {
				m, _ := v.(map[string]any)
				name, _ := m["fieldName"].(string)
				if original, ok := oldFields[strings.ToLower(name)]; ok {
					for _, flag := range []string{"isNeedEncrypt", "isSensitive"} {
						prior, _ := original[flag].(bool)
						current, _ := m[flag].(bool)
						if prior != current {
							add(fmt.Sprintf("body.objectType.fields[%d].%s", i, flag), "Encryption and sensitivity attributes cannot change after a field is created.")
						}
					}
				}
			}
			oldIndexes := map[string]any{}
			oldItems, _ := before["indexes"].([]any)
			for _, v := range oldItems {
				m, _ := v.(map[string]any)
				name, _ := m["indexName"].(string)
				oldIndexes[strings.ToLower(name)] = m["indexList"]
			}
			newItems, _ := after["indexes"].([]any)
			for i, v := range newItems {
				m, _ := v.(map[string]any)
				name, _ := m["indexName"].(string)
				if prior, ok := oldIndexes[strings.ToLower(name)]; ok {
					a, _ := json.Marshal(prior)
					b, _ := json.Marshal(m["indexList"])
					if !bytes.Equal(a, b) {
						add(fmt.Sprintf("body.objectType.indexes[%d].indexList", i), "Existing index definitions cannot be modified in place; the documented editor supports adding or deleting indexes.")
					}
				}
			}
		}
	}
	if e.FamilyID == "templates" && (e.ID == "create-application" || e.ID == "create-plan" || e.ID == "deploy") {
		components, _ := root["components"].([]any)
		seen := map[string]bool{}
		for i, v := range components {
			m, _ := v.(map[string]any)
			name, _ := m["name"].(string)
			if seen[name] {
				add(fmt.Sprintf("body.components[%d].name", i), "SAM component names must be unique.")
			}
			seen[name] = true
			props, _ := m["properties"].(map[string]any)
			if events, ok := props["events"].([]any); ok {
				names := map[string]bool{}
				for j, item := range events {
					event, _ := item.(map[string]any)
					name, _ := event["name"].(string)
					prefix := fmt.Sprintf("body.components[%d].properties.events[%d]", i, j)
					if names[name] {
						add(prefix+".name", "SAM event names must be unique within a component.")
					}
					names[name] = true
					if event["type"] == "cron" {
						properties, _ := event["properties"].(map[string]any)
						schedule, _ := properties["schedule"].(string)
						valid := regexp.MustCompile(`^@every [1-9][0-9]*[smh]$`).MatchString(schedule)
						if strings.HasPrefix(schedule, "cron(") && strings.HasSuffix(schedule, ")") {
							valid = validCron(strings.TrimSuffix(strings.TrimPrefix(schedule, "cron("), ")"))
						}
						if !valid {
							add(prefix+".properties.schedule", "Use the documented @every or five-field cron(...) schedule.")
						}
						input, _ := properties["input"].(string)
						if !json.Valid([]byte(input)) {
							add(prefix+".properties.input", "Cron event input must remain a string containing valid JSON.")
						}
					}
				}
			}
			if description, ok := props["description"].(string); ok {
				width := 0
				for _, unit := range utf16.Encode([]rune(description)) {
					if unit <= 128 {
						width++
					} else {
						width += 2
					}
				}
				if width > 64 {
					add(fmt.Sprintf("body.components[%d].properties.description", i), "SAM descriptions have a maximum vendor display width of 64 (non-ASCII UTF-16 units count twice).")
				}
			}
			if env, ok := props["env"].([]any); ok {
				for j, item := range env {
					entry, _ := item.(map[string]any)
					value, _ := entry["value"].(string)
					if len(utf16.Encode([]rune(value))) > 2048 {
						add(fmt.Sprintf("body.components[%d].properties.env[%d].value", i, j), "SAM environment values have a maximum of 2048 UTF-16 units.")
					}
				}
			}
			uri, _ := props["codeUri"].(string)
			if e.ID != "deploy" && (strings.Contains(uri, "://") || strings.HasPrefix(uri, "/") || !strings.Contains(uri, "/")) {
				add(fmt.Sprintf("body.components[%d].properties.codeUri", i), "Provide an already uploaded bucket/object.zip key, not a local path or URL.")
			}
		}
	}
	return issues
}

func validCron(expression string) bool {
	fields := strings.Fields(expression)
	if len(fields) != 5 {
		return false
	}
	bounds := [][2]int{{0, 59}, {0, 23}, {1, 31}, {1, 12}, {0, 6}}
	for i, field := range fields {
		upper := strings.ToUpper(field)
		if i == 3 {
			for j, n := range []string{"JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC"} {
				upper = strings.ReplaceAll(upper, n, strconv.Itoa(j+1))
			}
		}
		if i == 4 {
			for j, n := range []string{"SUN", "MON", "TUE", "WED", "THU", "FRI", "SAT"} {
				upper = strings.ReplaceAll(upper, n, strconv.Itoa(j))
			}
			if strings.Contains(upper, "/") {
				return false
			}
		}
		for _, part := range strings.Split(upper, ",") {
			steps := strings.Split(part, "/")
			if len(steps) > 2 {
				return false
			}
			if len(steps) == 2 {
				n, err := strconv.Atoi(steps[1])
				if err != nil || n <= 0 {
					return false
				}
			}
			if steps[0] == "*" {
				continue
			}
			r := strings.Split(steps[0], "-")
			if len(r) > 2 {
				return false
			}
			prior := -1
			for _, v := range r {
				n, err := strconv.Atoi(v)
				if err != nil || n < bounds[i][0] || n > bounds[i][1] || prior > n {
					return false
				}
				prior = n
			}
		}
	}
	return true
}

func validateSchemaModel(m map[string]any, path string, add func(string, string)) {
	name, _ := m["objectTypeName"].(string)
	reserved := " naturalbase_metadata objecttypeinfohelper t_data_upgrade_info t_index_schema t_nstore_config t_schema_negotiate_info t_metadata_schema t_nstore_permission t_system_config "
	if strings.HasPrefix(strings.ToLower(name), "sqlite_") || strings.Contains(reserved, " "+strings.ToLower(name)+" ") {
		add(path+".objectTypeName", "Use a non-reserved object type name; sqlite_ and system table names are forbidden.")
	}
	fields, _ := m["fields"].([]any)
	primary, encrypted, auto := 0, 0, 0
	seen := map[string]bool{}
	types := map[string]string{}
	sensitive := map[string]bool{}
	for i, v := range fields {
		f, _ := v.(map[string]any)
		n, _ := f["fieldName"].(string)
		typ, _ := f["fieldType"].(string)
		fp := fmt.Sprintf("%s.fields[%d]", path, i)
		if seen[strings.ToLower(n)] {
			add(fp+".fieldName", "Field names are case-insensitively unique.")
		}
		seen[strings.ToLower(n)] = true
		types[n] = typ
		reservedFields := " naturalbase_version naturalbase_deleted naturalbase_operationtype naturalbase_creator naturalbase_accesstime naturalbase_operationtime naturalbase_syncstatus naturalbase_changedfieldsbitmap naturalbase_lastmodifier cmin cmax xmin xmax ctid oid tableoid xc_node_id tablebucketid rowid "
		if strings.Contains(reservedFields, " "+strings.ToLower(n)+" ") {
			add(fp+".fieldName", "System field names are reserved.")
		}
		pk, _ := f["belongPrimaryKey"].(bool)
		enc, _ := f["isNeedEncrypt"].(bool)
		sens, _ := f["isSensitive"].(bool)
		nn, _ := f["notNull"].(bool)
		_, hasDefault := f["defaultValue"]
		sensitive[n] = sens
		if pk {
			primary++
			if strings.Contains(" ByteArray Text Date Double Float Boolean ", " "+typ+" ") {
				add(fp+".belongPrimaryKey", "This field type cannot be a primary key.")
			}
		}
		isAuto := typ == "IntAutoIncrement" || typ == "LongAutoIncrement"
		if isAuto {
			auto++
		}
		if enc || sens {
			encrypted++
			if pk || nn || hasDefault || enc && sens {
				add(fp+".isNeedEncrypt", "Encrypted or sensitive fields cannot be primary, non-null, defaulted, or both encrypted and sensitive.")
			}
		}
		if enc && strings.Contains(" ByteArray Text IntAutoIncrement LongAutoIncrement ", " "+typ+" ") {
			add(fp+".isNeedEncrypt", "This field type cannot be encrypted.")
		}
		if sens && !strings.Contains(" Byte Short Integer Long Float Double String Date ", " "+typ+" ") {
			add(fp+".isSensitive", "Unsupported sensitive-field type.")
		}
		if nn && (typ == "ByteArray" || typ == "Date") {
			add(fp+".notNull", "ByteArray and Date cannot be explicitly non-null.")
		}
		if nn && !pk && !isAuto && !hasDefault {
			add(fp+".defaultValue", "Non-null, non-primary fields require a default value.")
		}
		if hasDefault && (pk || enc || sens || isAuto || typ == "ByteArray" || typ == "Date") {
			add(fp+".defaultValue", "This field cannot have a default value.")
		}
		if def, ok := f["defaultValue"].(string); ok && typ == "Text" && utf8.RuneCountInString(def) > 200 {
			add(fp+".defaultValue", "Text defaults have a maximum of 200 characters.")
		}
	}
	if primary < 1 || primary > 5 {
		add(path+".fields", "Define one to five primary-key fields.")
	}
	if encrypted > 5 {
		add(path+".fields", "At most five encrypted and sensitive fields combined.")
	}
	if auto > 1 {
		add(path+".fields", "At most one auto-increment field is supported.")
	}
	indexes, _ := m["indexes"].([]any)
	seenIndex := map[string]bool{}
	for i, v := range indexes {
		index, _ := v.(map[string]any)
		name, _ := index["indexName"].(string)
		ip := fmt.Sprintf("%s.indexes[%d]", path, i)
		if seenIndex[strings.ToLower(name)] {
			add(ip+".indexName", "Index names must be case-insensitively unique.")
		}
		seenIndex[strings.ToLower(name)] = true
		list, _ := index["indexList"].([]any)
		for j, f := range list {
			field, _ := f.(map[string]any)
			name, _ := field["fieldName"].(string)
			if types[name] == "" || types[name] == "ByteArray" || types[name] == "Text" || sensitive[name] {
				add(fmt.Sprintf("%s.indexList[%d].fieldName", ip, j), "Index must reference an existing indexable non-sensitive field.")
			}
		}
	}
}
