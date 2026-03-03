// Základní JavaScript pro canarias.run MVP

document.addEventListener('DOMContentLoaded', () => {
    
    // Zástupná logika pro Language Switcher
    const langBtns = document.querySelectorAll('.lang-btn');
    langBtns.forEach(btn => {
        btn.addEventListener('click', (e) => {
            // Remove active class
            langBtns.forEach(b => b.classList.remove('active'));
            // Set active
            e.target.classList.add('active');
            
            // ToDo: Here we would trigger local storage or URL change for language
            console.log('Language switched to:', e.target.innerText);
        });
    });

    // Zástupná logika pro formulář
    const searchForm = document.querySelector('.search-box');
    if (searchForm) {
        searchForm.addEventListener('submit', (e) => {
            e.preventDefault();
            const searchInput = searchForm.querySelector('.search-input').value;
            const island = searchForm.querySelectorAll('.filter-select')[0].value;
            const type = searchForm.querySelectorAll('.filter-select')[1].value;
            const month = searchForm.querySelectorAll('.filter-select')[2].value;

            console.log('Vyhledávám:', { searchInput, island, type, month });
            // Zlaté tlačítko dostane malou animaci načítání
            const btnSubmit = searchForm.querySelector('.btn-search');
            const originalText = btnSubmit.innerHTML;
            btnSubmit.innerHTML = '<span class="icon-search">⏳</span>';
            
            setTimeout(() => {
                btnSubmit.innerHTML = originalText;
                alert('Tato MVP verze zatím nemá funkční backendové filtrování. Parametry: ' + JSON.stringify({ searchInput, island, type, month }));
            }, 800);
        });
    }

    // Quick islands click
    const islandBadges = document.querySelectorAll('.island-badge');
    islandBadges.forEach(badge => {
        badge.addEventListener('click', (e) => {
            const islandText = e.target.innerText.replace(/[\uE000-\uF8FF]|\uD83C[\uDF00-\uDFFF]|\uD83D[\uDC00-\uDDFF]/g, '').trim(); // Remove emoji
            console.log('Quick filter pre:', islandText);
        });
    });
});
