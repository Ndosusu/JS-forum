export function createProfile(parentElement) {
    const container = document.createElement("div");
    container.classList.add("container");

    const pageTitle = document.createElement("h2");
    pageTitle.textContent = "Profile";
    container.appendChild(pageTitle);


    const backButton = document.createElement("button");
    backButton.classList.add("backHome")
    backButton.textContent = "Retour à l'accueil";
    backButton.addEventListener("click", () => {
        loadPage("index");
    });

    container.appendChild(backButton);
    
    parentElement.appendChild(container);
}