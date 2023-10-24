const TODOS = []

function addTodo() {
    let newTodo = {};
    newTodo['value'] = 0;
}

function clearInputs() {
    let inputs = document.getElementsByTagName("input");
    for (let i = 0; i < inputs.length; i++) {
        inputs[i].value = '';
    }
}

function startCountUp(countUp) {
    let startButton = `todoCountUp${countUp}`;
    let elem = window[startButton];
    let interval = setInterval(() => {
        if (!elem.value) elem.value = 0;
        elem.innerHTML = formatTime(elem.value++);
    }, 1000);

    let stopButton = `template${countUp}`;
    let elm = window[stopButton].content.cloneNode(true).querySelector("button");
    elm.onclick = stopCounter.bind(elm, interval);
    let start = window[`start${countUp}`];
    start.replaceChild(elm, start.children[0]);

    console.log(window[`start${countUp}`]);
    // window[startButton] = elm;
}

function stopCounter(interval, button) {
    console.log("Stopping counter", button);
    clearInterval(interval);

}