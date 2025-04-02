import {parseRecipe} from './recipe.js'

const supportsHTML5 = ('content' in document.createElement('template'))
const endpoint = 'https://127.0.0.1:443/add-recipe';

document.addEventListener('DOMContentLoaded', () => {
    if (!supportsHTML5) {
        // TODO: tell user they need to use a modern browser
    }

    const form = document.getElementById("recipe_form");
    form.addEventListener('submit', (e) => {
        e.preventDefault();

        let recipe = parseRecipe();
        let json = JSON.stringify(recipe);
        let r = new Request(endpoint, {
            method: 'POST',
            body: json,
            headers: {
                'content-type': 'application/json'
            }
        });

        fetch(r)
            .then(r => {
                // TODO: handle it!
                if (!r.ok) throw new Error('invalid');
                console.log(r.text());
            })
            .catch(console.warn);
    });
});