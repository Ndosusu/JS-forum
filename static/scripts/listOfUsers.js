import { GetUserInfoFromCookie } from "./utils.js";

export function createListOfUsers() {
    const list = document.createElement("div");
    list.classList.add("listContainer");
    
    list.innerHTML = "<h1>Utilisateurs connectés</h1>";

    const me = GetUserInfoFromCookie()

    /*
    list.innerHTML = `
    <h1>Liste des Utilisateurs</h1>
    <ul id="usersList">
        <h2>Online</h2>
        <li class="online"><span class="status"></span> Bob</li>
        <li class="online"><span class="status"></span> Bobby</li>
        <h2>Offline</h2>
        <li class="offline"><span class="status"></span> Charles</li>
        <li class="offline"><span class="status"></span> Charly</li>
    </ul>
    `
    */

    const USERS_LIST = document.createElement("ul");
    USERS_LIST.id = "usersList";

    fetch("/api/whosConnected")
        .then(response => response.json())
        .then(users => {
            if (users.length <= 1) {
                const NO_USER = document.createElement("li");
                NO_USER.classList.add("no-click");
                NO_USER.innerHTML = "Personne n'est connecté :(";
                USERS_LIST.appendChild(NO_USER);
            } else {
                users.forEach(user => {
                    if (user.username !== me.username) {
                        const USER_ITEM = document.createElement("li");
                        USER_ITEM.classList.add("user");
                        USER_ITEM.innerHTML = user.username;
                        // TODO: Faire apparaître l'onglet de DMs avec cette personne
                        //USER_ITEM.onclick = () => {
                        // 
                        // };
                        USERS_LIST.appendChild(USER_ITEM);
                    } 
                });
            }

            list.appendChild(USERS_LIST);
            document.body.prepend(list);
        })
}