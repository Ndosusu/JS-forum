import { LoadPage } from "./utils.js";

export function createProfile(parentElement) {
    const container = document.createElement("div");
    container.classList.add("container");

    const pageTitle = document.createElement("h2");
    pageTitle.textContent = "Profile";
    container.appendChild(pageTitle);


    const backButton = document.createElement("button");
    backButton.textContent = "Déconnexion";
    backButton.addEventListener("click", async () => {
        const response = await fetch("api/logout")
        if (!response.ok) {
            alert("Error upon disconnect...")
            return
        } else {
            LoadPage("registration");
        }
    });

    container.appendChild(backButton);
    
    parentElement.appendChild(container);
}