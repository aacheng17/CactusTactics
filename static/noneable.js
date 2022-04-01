export class Noneable {
    constructor(element, showAtStart = false) {
        this.element = element;
        this.display = element.style.display || "block"; // any noneable that we define a display for needs to be defined in the ingame.html because of how the css is imported
        if (!showAtStart) {
            this.element.style.display = "none";
        }
    }

    show() {
        this.element.style.display = this.display;
    }

    hide() {
        this.element.style.display = "none";
    }
  }