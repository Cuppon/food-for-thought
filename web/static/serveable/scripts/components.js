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
// TODO: look into whether shadowRoot detection necessary (check adoptedStyleSheets in Document.prototype)
const shadowSheet = new CSSStyleSheet();
shadowSheet.replaceSync(cssRules);
const LABEL_TEXT = 'label-text';

class CustomBase extends HTMLElement {
    constructor() {
        super();
    }

    get value() {
        throw new Error("Child class must implement the 'value' property.")
    }
}

// TODO: look into ARIA attributes
const checkboxItemTemplate = document.createElement('template');
checkboxItemTemplate.innerHTML = `
    <div class="flex flex-col flex-1 text-center">
        <div class="flex-1">
            <p></p>
        </div>
        <div class="flex-1 px-4 ">
            <label class="block text-gray-50 text-xs leading-none"></label>
        </div>
        <div class="flex-1">
            <input type="checkbox">
        </div>
    </div>
`;

class CheckboxItem extends CustomBase {
    _shadowRoot;
    _input;

    constructor() {
        super();

        this._shadowRoot = this.attachShadow({mode:'closed'});
        let clone = checkboxItemTemplate.content.cloneNode(true);
        this._shadowRoot.append(clone);
        this._shadowRoot.adoptedStyleSheets = [shadowSheet];
    }

    get value() {
        return this._input.checked;
    }

    connectedCallback() {
        this._input = this._shadowRoot.querySelector('input');
        this._input.name = this.id;

        const labelText = this.getAttribute(LABEL_TEXT);
        const label = this._shadowRoot.querySelector('label');
        label.setAttribute('for', this.id);
        label.textContent = labelText;

        const emojiIcon = this.getAttribute('emoji-code');
        const p = this._shadowRoot.querySelector('p');
        p.textContent = emojiIcon;
    }
}

const containerGroupTemplate = document.createElement('template');
containerGroupTemplate.innerHTML = `
    <legend class="block text-gray-50"></legend>
    <div class="flex flex-row">
        <slot></slot>
    </div>
    <style>
        ::slotted(*) { display: flex; flex: 1 1; }
    </style>
`;

class ContainerGroup extends CustomBase {
    _shadowRoot;
    _legend;
    _toggleableElement;
    _children;

    constructor() {
        super();

        this._shadowRoot = this.attachShadow({mode:'closed'});
        let clone = containerGroupTemplate.content.cloneNode(true);
        this._shadowRoot.append(clone);
        this._shadowRoot.adoptedStyleSheets = [shadowSheet];
    }

    get value() {
        if (!this._children) throw new Error(`Can't get ${this.className} value; uninitialized`);

        let children = {};
        for(const c of this._children) {
            children[c.id] = c.value;
        }

        return children;
    }

    connectedCallback() {
        this._toggleableElement = this._shadowRoot.querySelector('div');

        this._legend = this._shadowRoot.querySelector('legend')
        this._legend.setAttribute('for', this.id)
        this._legend.textContent = this.getAttribute(LABEL_TEXT);

        this._children = this._shadowRoot.querySelector('slot').assignedElements();
    }

    toggleVisibility() {
        const show = this._toggleableElement.style.display === 'none';
        this._toggleableElement.style.display = show ? 'flex' : 'none';
    }
}

const textInputTemplate = document.createElement('template');
textInputTemplate.innerHTML = `
    <div class="flex items-center">
        <label class="text-gray-50"></label>
        <input type="text" required class="ml-2">
    </div>
`;

class TextInput extends CustomBase {
    _shadowRoot;
    _toggleableElement;

    constructor() {
        super();

        this._shadowRoot = this.attachShadow({mode: 'closed'});
        let clone = textInputTemplate.content.cloneNode(true);
        this._shadowRoot.append(clone);
        this._shadowRoot.adoptedStyleSheets = [shadowSheet];
    }

    get value() {
        if (!this._toggleableElement) throw new Error(`Can't get ${this.className} value; uninitialized`);

        return this._toggleableElement.value;
    }

    connectedCallback() {
        this._toggleableElement = this._shadowRoot.querySelector('input');
        this._toggleableElement.name = this.id;

        const label = this._shadowRoot.querySelector('label');
        label.setAttribute('for', this.id);
        label.textContent = this.getAttribute(LABEL_TEXT);
    }

    toggleVisibility() {
        const show = this._toggleableElement.style.display === 'none';
        this._toggleableElement.style.display = show ? 'flex' : 'none';
    }
}

// TODO: look into ARIA attributes
const togglerTemplate = document.createElement('template');
togglerTemplate.innerHTML = `
    <div class="flex mt-4">
        <button class="flex-shrink-0 rounded-full w-9 h-9 border-2 border-gray-600 text-sm text-gray-100 hover:bg-gray-600">+</button>
        <div class="w-full ps-2 pt-1">
            <slot></slot>
        </div>
    </div>
`;

class ComponentToggler extends HTMLElement {
    static observedAttributes = ['show-initially'];
    _shadowRoot;
    _children;
    _definedChildren = [];

    constructor() {
        super();

        this._shadowRoot = this.attachShadow({mode: 'closed'});
        let clone = togglerTemplate.content.cloneNode(true);
        this._shadowRoot.append(clone);
        this._shadowRoot.adoptedStyleSheets = [shadowSheet];
    }

    connectedCallback() {
        const button = this._shadowRoot.querySelector('button');
        button.addEventListener('click', this.#toggleButtonDisplay);
        button.addEventListener('click', this.#toggleSlotVisibility.bind(this));

        this._children = this._shadowRoot.querySelector('slot').assignedElements();

        // hide slots by default unless attribute exists to specify otherwise
        const showInitially = this.hasAttribute('show-initially');
        if (!showInitially) {
            button.textContent = '+';

            this.#toggleSlotVisibility();
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

    /* TODO: need to make sure to check whether browser supports slots
     * - double check on toggling ARIA attributes when visibility toggled
     */
    async #toggleSlotVisibility() {
        for(const c of this._children) {
            const tag = c.tagName.toLowerCase();
            if (!this._definedChildren.includes(tag)) {
                this._definedChildren.push(tag);

                await customElements.whenDefined(tag);
            }

            if (typeof c.toggleVisibility === 'function') {
                c.toggleVisibility();
            } else {
                const errMsg = "Missing toggleVisibility() on " + c.name;
                throw new Error(errMsg);
            }
        }
    }

    toggleVisibility() {
        const show = this.style.display === 'none';
        this.style.display = show ? 'flex' : 'none';
    }
}

customElements.define('checkbox-item', CheckboxItem);
customElements.define('text-input', TextInput);
customElements.define('container-group', ContainerGroup);
customElements.define('component-toggler', ComponentToggler);
