export class Noneable {
    constructor(element, showAtStart = false) {
        this.element = element;
        this.display = element.style.display || "block";
        if (!showAtStart) {
            this.element.style.display = "none";
        }
    }

    show() {
        console.log(this.display);
        this.element.style.display = this.display;
    }

    hide() {
        this.element.style.display = "none";
    }
  }