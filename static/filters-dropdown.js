const button = document.querySelector("button#show-checkboxes")
const container = document.querySelector("#filter-container")
const allEffectsCheckbox = container.querySelector("#all-effects")
const allPhonesCheckbox = container.querySelector("#all-phones")
const allCheckboxes = container.querySelectorAll(`label>input:not([name=''])`)

function toggleFiltersVisibility(e) {
    container.classList.toggle("open")
    button.textContent = container.classList.contains("open") ? "Hide Filters" : "Show Filters"
}

function toggleAllCheckboxes(e) {
    let toggleAllCheckbox = e.target
    let checkboxes = container.querySelectorAll(`label:has(#${toggleAllCheckbox.id}) ~ label>input`)
    for (let checkbox of checkboxes)
        checkbox.checked = toggleAllCheckbox.checked
}

function toggleTheAllCheckbox(e) {
    let toggleAllCheckbox = e.target.parentElement.parentElement.firstChild.firstChild
    let checkboxes = container.querySelectorAll(`label:has(#${toggleAllCheckbox.id}) ~ label>input`)

    let all = true
    for (let checkbox of checkboxes) {
        if (!checkbox.checked) {
            all = false
            break
        }
    }
    toggleAllCheckbox.checked = all
}

button.addEventListener("click", toggleFiltersVisibility)
allEffectsCheckbox.addEventListener("change", toggleAllCheckboxes)
allPhonesCheckbox.addEventListener("change", toggleAllCheckboxes)
container.addEventListener("change", toggleTheAllCheckbox)
window.addEventListener("load", function (e) {
    let phoneCheckboxes = container.querySelectorAll(`label:has(#all-phones) ~ label>input`)
    let effectCheckboxes = container.querySelectorAll(`label:has(#all-effects) ~ label>input`)

    let all = true
    for (let checkbox of phoneCheckboxes) {
        if (!checkbox.checked) {
            all = false
            break
        }
    }
    allPhonesCheckbox.checked = all

    all = true
    for (let checkbox of effectCheckboxes) {
        if (!checkbox.checked) {
            all = false
            break
        }
    }
    allEffectsCheckbox.checked = all
})