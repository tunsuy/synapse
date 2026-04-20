# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| v0.2.x  | ✅ Active  |
| v0.1.x  | ⚠️ Security fixes only |
| < v0.1  | ❌ Not supported |

## Reporting a Vulnerability

We take the security of Synapse seriously. If you discover a security vulnerability, please follow the responsible disclosure process below.

### ⚠️ Do NOT Create a Public Issue

Security vulnerabilities should **NOT** be reported through public GitHub issues. Public disclosure gives attackers a head start.

### How to Report

1. **Email**: Send a detailed report to **957126743@qq.com** with the subject line `[SECURITY] Synapse Vulnerability Report`
2. **Include**:
   - Description of the vulnerability
   - Steps to reproduce
   - Affected versions
   - Potential impact assessment
   - Suggested fix (if any)

### Response Timeline

| Stage | Timeline |
|-------|----------|
| Acknowledgment | Within **48 hours** |
| Initial assessment | Within **5 business days** |
| Fix development | Based on severity (see below) |
| Public disclosure | After fix is released |

### Severity Levels

| Severity | Description | Fix Timeline |
|----------|-------------|-------------|
| **Critical** | Remote code execution, data exfiltration, authentication bypass | Within **7 days** |
| **High** | Plugin sandbox escape, path traversal, privilege escalation | Within **14 days** |
| **Medium** | Information disclosure, denial of service | Within **30 days** |
| **Low** | Minor issues with limited impact | Next release |

## Security Design Principles

Synapse follows these security principles by design:

### 🔒 Data Sovereignty

- **Local-first by default**: The default Local Store keeps all data on the user's machine
- **No telemetry**: Synapse does not collect or transmit any user data
- **User-controlled storage**: Users choose where their knowledge lives (local / GitHub / S3 / etc.)

### 🧩 Plugin Isolation

- **Subprocess sandboxing**: External plugins run in separate processes, communicating via JSON-RPC over stdin/stdout
- **No direct filesystem access**: Plugins can only access data through the defined extension point interfaces
- **Path traversal protection**: Store implementations validate paths to prevent directory escape
- **Plugin name anti-spoofing**: Reserved name set + regex detection prevent malicious plugins from impersonating official ones

### ✅ Integrity Verification

- **Checksum validation**: Plugin installation verifies integrity via checksums
- **Version pinning**: Dependencies are version-locked for deterministic builds
- **Configuration validation**: `synapse check` verifies all config references are valid and safe

### 🔑 Credential Management

- **Environment variable references**: Config files use `${ENV_VAR}` syntax to avoid hardcoding secrets
- **No credential storage**: Synapse never stores tokens or passwords in plain text
- **Minimal permissions**: GitHub Store uses the minimum required token scope

## Security Best Practices for Users

1. **Use environment variables** for sensitive config values (GitHub tokens, API keys)
2. **Review third-party plugins** before installation
3. **Keep Synapse updated** to receive security patches
4. **Use Local Store** for sensitive knowledge that should not leave your machine
5. **Set restrictive file permissions** on `~/.synapse/config.yaml` (e.g., `chmod 600`)

## Acknowledgments

We appreciate the security research community's efforts in helping keep Synapse secure. Reporters of valid vulnerabilities will be acknowledged in our release notes (unless they prefer to remain anonymous).
