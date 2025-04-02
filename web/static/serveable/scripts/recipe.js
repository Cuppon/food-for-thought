// TODO: update to contain full recipe structure
const recipeStructure = {
    english_name: "",
    native_name: "",
    url: ""
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