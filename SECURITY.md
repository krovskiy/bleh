# security vulnerabilities


## 1. **no session expiration**


**vulnerability:**
sessions never expire. A stolen or compromised session token remains valid indefinitely.

**how to Exploit:**
1. steal a session cookie (via XSS, network sniffing, or compromised storage)
2. attacker has unlimited time to exploit stolen credentials

**impact:** privilege escalation

**fix:** add expiration timestamps:

```go
// previously

// nothing...
```

```go
// after
type SessionData struct {
    UserID    int
    ExpiresAt time.Time
}
var Sessions = map[string]SessionData{}

...

sessionData, exists := services.Sessions[session.Value]
if !exists || time.Now().After(sessionData.ExpiresAt) {
    http.Redirect(w, r, "/login/", http.StatusFound)
    return
}
```