// TODO: update to contain full recipe structure, and implement nested object logic, e.g. attribution { url: {}}
const recipeStructure = {
    english_name: "",
    native_name: "",
    attribution: "",
    url: "",
    emojis: ""
}

export function parseRecipe() {
    return Object.keys(recipeStructure).reduce((recipe, key) => {
        const e = document.getElementById(key);
        if (e) {
            recipe[key] = e.value;
        }
        return recipe;
    }, {});
}