const supportsHTML5 = ('content' in document.createElement('template'))

document.addEventListener('DOMContentLoaded', () => {
    if(!supportsHTML5) {
        // TODO: tell user they need to use a modern browser
    }

    // TODO: update names
    const form = document.getElementById("myform");
    const btn = document.getElementById("mybtn");

    btn.addEventListener('click', () => {
        // TODO: get all elements within the form that have a shared/common class, indicating
        // they are a component that has input data associated with them
        const textInput = form.querySelector('text-input');
        if(textInput) {
            let inputVal = textInput.value;
            console.log("and the value is ", inputVal);
        }
    })
})