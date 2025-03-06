export function createListOfUsers () {
    
    const list = document.createElement("div");
    list.classList.add("listContainer");
    list.innerHTML = `
    <h1>Liste des Utilisateurs</h1>
<ul>
    <h2>Online</h2>
    <li class="online"><span class="status"></span> Bob</li>
    <li class="online"><span class="status"></span> Bobby</li>
    <h2>Offline</h2>
    <li class="offline"><span class="status"></span> Charles</li>
    <li class="offline"><span class="status"></span> Charly</li>
</ul>
    `

    document.body.prepend(list);
}