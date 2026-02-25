const mainCont = document.getElementById("container");
const submitBtn = document.getElementById("submit");
const title = document.getElementById("title");
const completed = document.getElementById("completed");


function displayCond(preProc) {
    if (preProc == true){
        return "Done!"
    } else {
        return "Not done!"
    }
}

function addTaskButtons(task, taskDiv) {
  ['Edit', 'Remove'].forEach(label => {
    const btn = document.createElement("button");
    btn.id = `${label.toLowerCase().slice(0, 4)}-${task.id}`;
    btn.type = "button";
    btn.textContent = label;
    taskDiv.append(btn);
  });
}

function formatTime(a) {
    var date = new Date(a * 1000).toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit', second: '2-digit' });
    return date;        
}

document.addEventListener("DOMContentLoaded", async(e) => {
    e.preventDefault();
    try {
        const response = await fetch('/tasks', {
            method: 'get',
            headers: {
                'Content-Type': 'application/json'
            }
        });
        
        if (response.ok) {
            const tasks = await response.json();
            tasks.forEach(task => {
                const taskDiv = document.createElement("div");
                const createTaskContent = document.createElement("div");
                const createTaskTime = document.createElement("div");
                const createTaskComplete = document.createElement("div");

                taskDiv.id = `task-${task.id}`;
                createTaskContent.id = `taskContent-${task.id}`;
                createTaskTime.id = `time-${task.id}`;
                createTaskComplete.id = `complete-${task.id}`;

                createTaskContent.textContent = `${task.title}`;
                
                createTaskTime.textContent = `${formatTime(task.time)}`;

                createTaskComplete.textContent = `${task.completed}`;
                
                taskDiv.append(createTaskContent, createTaskTime, createTaskComplete);
                addTaskButtons(task, taskDiv);
                mainCont.append(taskDiv);
                
            });   
        }
    } catch(err){
        console.log(`Error: ${err}`);
    }
});

submitBtn.addEventListener("click", async(e) => {
    e.preventDefault();
    try {
        const response = await fetch('/tasks',{
            method: 'post',
            headers: {
                'Content-Type': 'application/json'
            }, //method, header, body
            body: JSON.stringify({ // remember to add stringify!
                title: title.value,
                completed: completed.checked,
            }),
        });
        if (response.ok){
            const newTask = await response.json();
            console.log("This task was created succesfully: ", newTask);
        } // check if it is okay
    } catch(err){
        console.error(`Error: ${err}`)
    } //display error otherwise
});

mainCont.addEventListener("click", async(e) => {
    elements = ['taskContent-', 'time-', 'complete-', 'edit-', 'remo-', 'temp-input'];
    
    const lastTask = e.target.id; //last clicked ID
    const tempInp = document.getElementById(`temp-input`);

    if (e.target.id.startsWith('remo-')) {
        const taskID = e.target.id.match(/remo-(\d+)/)?.[1];
        try {
            const response = await fetch(`/tasks/${taskID}`, {
                method: 'DELETE',
                headers: { 
                    'Content-Type' : 'application/json'
                }
            });
            if (response.ok){
                document.getElementById(`task-${taskID}`).remove();
            }

        } catch(err){
            console.error(`Error: ${err}`);
        }
        const a = document.getElementById('temp-input');
        if (a) {
        a.remove();
        }
        return;
    }

    if (e.target.id.startsWith('edit-')) {
        const taskID = e.target.id.match(/edit-(\d+)/)?.[1];
        const a = document.getElementById(`temp-input`);
        if (a == null){
            console.log("Error: no input!");  
            return;  
        }
        
        try {
            const response = await fetch(`/tasks/${taskID}`, {
                method: 'PUT',
                headers: { 
                    'Content-Type' : 'application/json'
                },
                body: JSON.stringify({
                    title: a.value
                })
            });
            if (response.ok){
                const updatedTask = await response.json();
                document.getElementById(`taskContent-${taskID}`).textContent = updatedTask.title;
                a.remove();
            }

        } catch(err){
            console.error(`Error: ${err}`);
        }
        return;
    }

    if (elements.some(prefix => e.target.id.startsWith(prefix))) {
        if (!tempInp){
            const a = document.createElement("input");
            a.setAttribute("id", `temp-input`);
            a.onclick = (event) => event.stopPropagation();
            mainCont.append(a);
            a.focus();
        }
    }

    if (!elements.some(prefix => e.target.id.startsWith(prefix))){
        console.log("Not found task!");
        return;
    }
    
});



document.addEventListener("click", (e) => {
    const elements = ['taskContent-', 'time-', 'complete-', 'edit-', 'remo-', 'temp-input-'];
    const isTaskElement = elements.some(prefix => e.target.id.startsWith(prefix));
    
    if (!isTaskElement) {
        document.querySelectorAll('[id^="temp-input"]').forEach(input => input.remove());
    }
});