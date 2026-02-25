submitBtn = document.getElementById("submit")
mainDiv = document.getElementById("container")

submitBtn.addEventListener("click", async(e) => {
    e.preventDefault();
    try {
        const response = await fetch('/api/login',{
            method: 'post',
            headers: {
                'Content-Type': 'application/json'
            }, 
            body: JSON.stringify({ 
                username: username.value,
                password: password.value,
            }),
        });
        if (response.ok){
            a = document.createElement('div');
            b = document.createElement('h1');
            b.textContent = 'succesfully sent!';
            a.append(b);
            mainDiv.append(a);
            window.location.href = "/";
        } 
    } catch(err){
        console.error(`Error: ${err}`)
    } 
});