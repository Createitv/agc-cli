package agcapi

import (
 "context"
 "io"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"

 "github.com/Createitv/agc-cli/pkg/domain"
)

func TestPrefetchJSONIsApplicationDataAndNeverReceivesDeveloperAuthorization(t *testing.T){
 data:=`{"code":17,"ret":{"code":23},"resources":["image.png"]}`
 srv:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.Method!="GET"||r.URL.Query().Get("appId")!="offline-app"||r.URL.Query().Get("params")!="{\"hello\":\"world\"}"{t.Error("prefetch request changed")}
  if r.Header.Get("Authorization")!=""||r.Header.Get("Cookie")!=""{t.Error("developer credentials reached callback")}
  io.WriteString(w,data)
 }));defer srv.Close()
 e,_:=domain.EndpointByID("prefetch","fetch-periodic")
 result,err:=InvokeEndpoint(context.Background(),nil,InvokeRequest{Endpoint:e,Params:map[string]string{"callbackUrl":srv.URL},Query:map[string]string{"appId":"offline-app","params":`{"hello":"world"}`,"token":"private-user-token"},AccessToken:"must-not-forward"})
 if err!=nil||string(result.RawBody)!=data{t.Fatalf("application data changed: %v %s",err,result.RawBody)}
 if strings.Contains(result.URL,"private-user-token"){t.Fatal("prefetch user token exposed in feedback URL")}
}

func TestPrefetchInstallRejectsResponseBeyondTwoMiB(t *testing.T){
 srv:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){io.WriteString(w,strings.Repeat("x",(2<<20)+1))}));defer srv.Close()
 e,_:=domain.EndpointByID("prefetch","fetch-install")
 _,err:=InvokeEndpoint(context.Background(),nil,InvokeRequest{Endpoint:e,Params:map[string]string{"callbackUrl":srv.URL},Query:map[string]string{"appId":"offline-app"}})
 if err==nil{t.Fatal("oversized install prefetch accepted")}
}
