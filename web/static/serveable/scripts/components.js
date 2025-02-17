// TODO: have this pull from specific component's css rules (css file). later, in component
// constructor, dont create a new stylesheet, rather assign adopted stylesheets to this one
const cssRules = (() => {
    const s = document.getElementById('layoutSheet').sheet;
    let cssTxt = '';
    for (let i = 0; i < s.cssRules.length; i++) {
        cssTxt += s.cssRules[i].cssText + ' ';
    }
    return cssTxt;
})();
const shadowSheet = new CSSStyleSheet();
shadowSheet.replaceSync(cssRules);

class CustomBase extends HTMLElement {
    constructor() {
        super();
    }

    get value() {
        throw new Error("Child class must implement the 'value' property.")
    }
}

const textInputTemplate = document.createElement('template');
textInputTemplate.innerHTML = `
    <div class="flex items-center space-x-2">
        <label class="text-gray-50"></label>
        <input type="text" required class="ml-2 w-15">
    </div>
`;

class TextInput extends CustomBase {
    _shadowRoot;
    _input;

    constructor() {
        super();

        this._shadowRoot = this.attachShadow({mode: 'closed'});
        let clone = textInputTemplate.content.cloneNode(true);
        this._shadowRoot.append(clone);
        this._shadowRoot.adoptedStyleSheets = [shadowSheet];
    }

    get value() {
        return this._input.value;
    }

    connectedCallback() {
        const inputName = this.getAttribute('input-name');
        this._input = this._shadowRoot.querySelector('input');
        this._input.id = inputName;
        this._input.name = inputName;

        const labelText = this.getAttribute('label-text');
        const label = this._shadowRoot.querySelector('label');
        label.setAttribute('for', inputName);
        label.textContent = labelText;
    }

    toggleVisibility() {
        const show = this._input.style.display === 'none';
        this._input.style.display = show ? 'flex' : 'none';
    }
}

const togglerTemplate = document.createElement('template');
togglerTemplate.innerHTML = `
    <div class="flex flex-row items-center ml-4">
        <button class="rounded-full w-12 h-12 text-2xl focus:outline-none">+</button>
        <div id="content" class="ml-4 flex-1 items-center space-x-4">
            <slot></slot>
        </div>
    </div>
`;

class ComponentToggler extends HTMLElement {
    static observedAttributes = ['show-initially'];
    _shadowRoot;

    constructor() {
        super();

        this._shadowRoot = this.attachShadow({mode: 'closed'});
        let clone = togglerTemplate.content.cloneNode(true);
        this._shadowRoot.append(clone);
        this._shadowRoot.adoptedStyleSheets = [shadowSheet];
    }

    connectedCallback() {
        const buttonID = this.getAttribute('id');
        const button = this._shadowRoot.querySelector('button');
        button.id = buttonID;

        button.addEventListener('click', this.#toggleButtonDisplay);
        button.addEventListener('click', this.#toggleSlotVisibility.bind(this));

        // hide slots by default unless attribute exists to specify otherwise
        const showInitially = this.hasAttribute('show-initially');
        if (!showInitially) {
            button.textContent = '+';

            this.#toggleSlotVisibility(null);
        } else {
            button.textContent = '-';
        }
    }

    disconnectedCallback() {
        const button = this._shadowRoot.querySelector('button');

        if (button) {
            button.removeEventListener('click', this.#toggleButtonDisplay);
            button.removeEventListener('click', this.#toggleSlotVisibility);
        }
    }

    #toggleButtonDisplay(ev) {
        const show = ev.target.textContent === '+';
        ev.target.textContent = show ? '-' : '+'
    }

    #toggleSlotVisibility(ev) {
        // ! - due to lifecycle and DOM order, moving this into connectedCallback(), would
        // result in 'children' being undefined
        const children = this._shadowRoot.querySelector('slot').assignedElements();

        children.forEach((e) => {
            if (typeof e.toggleVisibility === 'function') {
                e.toggleVisibility();
            } else {
                const errMsg = "Missing toggleVisibility() on " + e.name;
                throw new Error(errMsg);
            }
        })
    }
}

customElements.define('text-input', TextInput);
customElements.define('component-toggler', ComponentToggler);
