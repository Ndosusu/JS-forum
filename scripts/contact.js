export function createContact(parentElement){
    const container = document.createElement("div");
    container.classList.add("container");

    const pageTitle = document.createElement("h2");
    pageTitle.textContent = "Contact";
    container.appendChild(pageTitle);

    const secContainer = document.createElement("div");
    secContainer.classList.add("secContainer");
    container.appendChild(secContainer);

    const listOfUsers = document.createElement("ol")
    listOfUsers.classList.add("listUsers");
    listOfUsers.textContent = "Liste des utilisateurs";
    secContainer.appendChild(listOfUsers);

    const messageContainer = document.createElement("div");
    messageContainer.classList.add("messageContainer");
    messageContainer.textContent = "Your messages";
    secContainer.appendChild(messageContainer);

    const yourMessage = document.createElement("div");
    yourMessage.classList.add("message");
    yourMessage.textContent = "message";
    messageContainer.appendChild(yourMessage);

    const backButton = document.createElement("button");
    backButton.classList.add("backHome")
    backButton.textContent = "Retour à l'accueil";
    backButton.addEventListener("click", () => {
        loadPage("index");
    });

    container.appendChild(backButton);
    
    parentElement.appendChild(container);
}