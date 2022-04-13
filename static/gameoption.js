export class GameOption {
    constructor(div, id, text, max, min, step) {
        this.div = div;
        this.max = max;
        this.min = min;
        this.step = step;

        this.text = document.createElement("a");
        this.text.innerText = text;
        this.div.appendChild(this.text);

        this.element = document.createElement("input");
        this.element.classList.add("game-option");
        this.element.setAttribute("type", "number");
        this.element.setAttribute("onkeydown", "return false");
        this.element.setAttribute("style", "caret-color: transparent");
        this.element.setAttribute("id", id);
        this.div.appendChild(this.element);

        this.left = document.createElement("button");
        this.left.innerText = "<";
        this.left.onclick = this.element.value -= 
        this.div.appendChild(this.left);

        this.right = document.createElement("button");
        this.right.innerText = ">";
        this.div.appendChild(this.right);
    }

    getValue() {
        return parseInt(this.element.value);
    }

    setValue(value) {
        this.element.value = parseInt(value);
    }

    decrement() {
        this.element.value = parseInt(this.element.value) - this.step;
        if (this.element.value < this.min) {
            this.element.value = this.min;
        }
    }
    
    increment() {
        this.element.value = parseInt(this.element.value) + this.step;
        if (this.element.value > this.max) {
            this.element.value = this.max;
        }
    }
}