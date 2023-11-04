const currentUrl = window.location.href;

const games = ["idiotmouth", "fakeout", "standoff", "aaranagrams"];

for (const game of games) {
    document.getElementById(`link-${game}`).href = `${currentUrl}${game}`;
}