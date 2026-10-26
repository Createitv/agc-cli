"""Select a signing certificate by CSR public key, independently of chain order."""
import subprocess
import hashlib
import re
import ssl


def certificate_fingerprint_for_csr(certificate_path, csr_path):
    """证书链无序；只对与本次 CSR 公钥一致的实体证书计算指纹。"""
    key = subprocess.run(["openssl", "req", "-in", str(csr_path), "-pubkey", "-noout"], capture_output=True, timeout=20)
    if key.returncode or not key.stdout.strip():
        raise RuntimeError("CSR public key extraction failed")
    content = certificate_path.read_bytes()
    pattern = rb"-----BEGIN CERTIFICATE-----.*?-----END CERTIFICATE-----"
    blocks = re.findall(pattern, content, re.DOTALL)
    if not blocks:
        for command in (["openssl", "pkcs7", "-inform", "DER", "-in", str(certificate_path), "-print_certs"], ["openssl", "x509", "-inform", "DER", "-in", str(certificate_path), "-outform", "PEM"]):
            decoded = subprocess.run(command, capture_output=True, timeout=20)
            if decoded.returncode == 0:
                blocks = re.findall(pattern, decoded.stdout, re.DOTALL)
                if blocks: break
    matches = []
    for block in blocks:
        public = subprocess.run(["openssl", "x509", "-pubkey", "-noout"], input=block, capture_output=True, timeout=20)
        if public.returncode:
            raise RuntimeError("Malformed certificate in downloaded chain")
        if public.stdout.strip() == key.stdout.strip():
            matches.append(hashlib.sha256(ssl.PEM_cert_to_DER_cert(block.decode("ascii"))).hexdigest().upper())
    if len(matches) != 1:
        raise RuntimeError("Certificate chain must contain exactly one CSR-matching signing certificate")
    return matches[0]

