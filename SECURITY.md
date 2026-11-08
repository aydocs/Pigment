# Security policy

## Scope

pigment is a dependency-free Go terminal styling library. It does not run a
server, persist user data, authenticate users, or handle credentials.

## Reporting a vulnerability

Please do not disclose a suspected vulnerability in a public issue. Contact
the repository maintainers privately through the security contact configured
for the GitHub repository and include:

- the affected version or commit;
- a minimal reproduction;
- the security impact; and
- a proposed mitigation, if known.

Terminal escape-sequence injection caused by intentionally passing untrusted
text or URLs to terminal-output APIs should be reported with the exact input
and terminal emulator involved.

## Safe-use note

Applications should treat terminal text and URLs as untrusted when they come
from users or remote systems. `Link` protects control bytes in the URL, but it
does not sanitize arbitrary terminal control sequences in the displayed text.
