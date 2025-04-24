import { LoadPage, GetUserInfoFromCookie, IsUserInfoValid } from "./utils.js";

// Si cookie valide & serveur reconnait l'utilisateur, charge la page principale. Sinon, charge la page de connexion / inscription
document.addEventListener("DOMContentLoaded", async () => {
    let whereToGo = "";
    
    if (!IsUserInfoValid()) {
        whereToGo = "registration";
    } else {
        const response = await fetch("/api/imConnected")
        if (response.ok) {
            /*
            let answer = await response.text()
            answer = answer.slice(0, answer.length-1)
            console.log(`Server says \"${answer}\"`);
            */
            whereToGo = "index"
        } else {
            console.warn("Your cookie is valid, but the server says you're not connected. Please log in again.");
            document.cookie = ""
        }
    }

    LoadPage(whereToGo);
});
