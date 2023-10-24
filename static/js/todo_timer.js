function formatTime(ms) {
    // let timeString = [`${ms % 1000}ms`];
    let timeString = [];
    ms = Math.floor(ms / 1000);

    if (ms >= 0) {
        timeString.push(`${ms % 60}s`);
        ms = Math.floor(ms / 60);
    }

    if (ms > 0) {
        timeString.push(`${ms % 60}m`);
        ms = Math.floor(ms / 60);
    }

    if (ms > 0) {
        timeString.push(`${ms}h`)
    }

    return timeString.reverse().join("");
}

customElements.define('todo-timer',

    class extends HTMLElement {
        constructor() {
            super();

            let elapsed = this.querySelector("[slot='todo-elapsed']");
            elapsed.innerText = formatTime(parseInt(elapsed.getAttribute("value")));
        }
    })

customElements.define('timer-control',

    class extends HTMLElement {
        timer;
        timerID;
        todoID;
        constructor() {
            self = super();
            this.todoID = `Todo${this.id}`;
        }

        connectedCallback() {
            const btn = this.children[0];

            btn.addEventListener("click", async () => {
                if (btn.innerText == "Start") {
                    btn.innerText = "Stop";
                    this.timer = await this.waitForElement();
                    this.startTimer();
                } else {
                    btn.innerText = "Start";
                    this.stopTimer();
                }
            });
        }

        waitForElement() {
            return new Promise(resolve => {
                let observer = new MutationObserver((mutations) => {

                    mutations.forEach((mutation) => {
                        if (mutation.addedNodes.length == 0) return

                        for (let i = 0; i < mutation.addedNodes.length; i++) {
                            // do things to your newly added nodes here
                            let node = mutation.addedNodes[i]
                            if (node.id == this.todoID) {
                                observer.disconnect();
                                resolve(node.querySelector("[slot='todo-elapsed']"));
                            }
                        }
                    })
                })

                observer.observe(document.body, {
                    childList: true,
                    subtree: true,
                    attributes: false,
                    characterData: false
                })
            });
        }



        startTimer() {
            let startValue = parseInt(this.timer.getAttribute("value"));
            let start = Date.now() - startValue;
            this.timerID = setInterval(() => {
                const elapsed = Date.now() - start;
                this.timer.innerText = formatTime(elapsed);
            }, 16);
        }

        stopTimer() {
            clearInterval(this.timerID);
        }
    });


customElements.define('todo-control',
    class extends HTMLElement {
        constructor() {
            self = super();
            this.todoID = `Todo${this.id}`;
        }

        connectedCallback() {
            const btn = this.children[0];

            btn.addEventListener("click", async () => {
                if (btn.innerText == "Edit") {
                    btn.innerText = "Done";
                } else {
                    btn.innerText = "Edit"
                }
            });
        }
    })