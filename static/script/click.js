const menuItems = ['all', 'stripe', 'netflix', 'uber'];

const onMenuClick = ev => {
  const target = ev.target;

  console.log('🚀 ~ file: click.js:6 ~ onMenuClick ~ target:', target);

  if (target.tagName === 'INPUT' && menuItems.includes(target.value.toLowerCase())) {
    const prevActive = document.querySelector('input.active');

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

    //  htmx event listener
    document.body.addEventListener('htmx:configRequest', function (evt) {
      const isHeadless = document.getElementById('headless-scrapping')?.checked || false;
      evt.detail.parameters['headless'] = isHeadless;
    });
  }, 1500);
})();
