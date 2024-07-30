const menuItems = ['all', 'stripe', 'netflix', 'uber'];

const onMenuClick = ev => {
  const target = ev.target;

  console.log('🚀 ~ file: click.js:6 ~ onMenuClick ~ target:', target);

  // const spinner = document.getElementById('spinner')
  // if(spinner && spinner.style.display !== 'none') {
  //   return
  // }

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
    // const menuContainer = document.getElementById('scrapper-menu');

    // menuContainer.addEventListener('click', onMenuClick);

    //  htmx (request start) event listener
    document.body.addEventListener('htmx:configRequest', function (ev) {
      const isHeadless = document.getElementById('headless-scrapping')?.checked || false;
      ev.detail.parameters['headless'] = isHeadless;
      onMenuClick(ev);
    });
  }, 1500);
})();
