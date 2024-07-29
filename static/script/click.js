const menuItems = ['all', 'stripe', 'netflix', 'uber'];

const onMenuClick = ev => {
  const target = ev.target;

  if (target.tagName === 'BUTTON' && menuItems.includes(target.innerText.toLowerCase())) {
    const prevActive = document.querySelector('button.active');

    if (prevActive) {
      prevActive.classList.remove('active');
    }
    // new active btn
    target.classList.add('active');
  }
};

(() => {
  setTimeout(() => {
    const menuContainer = document.getElementById('scrapper-menu');

    menuContainer.addEventListener('click', onMenuClick);
  }, 1500);
})();
