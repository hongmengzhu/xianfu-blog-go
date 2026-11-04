/* AItool-Hub clone — interactions */
(function () {
  'use strict';

  const html = document.documentElement;

  function applyTheme(theme) {
    theme = theme === 'light' ? 'light' : 'dark';
    html.dataset.theme = theme;
    html.dataset.bsTheme = theme;
    html.classList.toggle('dark', theme === 'dark');
    try { localStorage.setItem('fby-theme', theme); } catch (e) {}
    document.querySelectorAll('[data-theme-btn]').forEach(function (btn) {
      btn.classList.toggle('active', btn.dataset.themeBtn === theme);
      btn.setAttribute('aria-pressed', btn.dataset.themeBtn === theme ? 'true' : 'false');
    });
    document.querySelectorAll('[data-theme-icon]').forEach(function (el) {
      el.setAttribute('href', '/assets/static/nav-001/icons.svg#' + (theme === 'light' ? 'sun' : 'moon'));
    });
  }

  function init() {
    const stored = (function () { try { return localStorage.getItem('fby-theme'); } catch (e) { return null; } })();
    applyTheme(stored || html.dataset.theme || 'dark');
    document.querySelectorAll('[data-theme-btn]').forEach(function (btn) {
      btn.addEventListener('click', function () { applyTheme(btn.dataset.themeBtn); });
    });

    /* Mobile sidebar */
    const body = document.body;
    const sidebar = document.getElementById('appSidebar');
    function closeSidebar() { body.classList.remove('sidebar-open'); if (sidebar) sidebar.setAttribute('aria-hidden', 'true'); }
    function openSidebar() { body.classList.add('sidebar-open'); if (sidebar) sidebar.setAttribute('aria-hidden', 'false'); }
    document.querySelectorAll('[data-sidebar-open]').forEach(function (b) { b.addEventListener('click', openSidebar); });
    document.querySelectorAll('[data-sidebar-close]').forEach(function (b) { b.addEventListener('click', closeSidebar); });
    const overlay = document.querySelector('.sidebar-overlay');
    if (overlay) overlay.addEventListener('click', closeSidebar);
    window.addEventListener('keydown', function (e) { if (e.key === 'Escape') closeSidebar(); });
    if (sidebar) sidebar.querySelectorAll('a').forEach(function (a) { a.addEventListener('click', function () { if (window.innerWidth < 1024) closeSidebar(); }); });

    /* Mobile search toggle */
    const mobileSearchToggle = document.querySelector('[data-mobile-search]');
    const topSearch = document.querySelector('.topbar-search');
    if (mobileSearchToggle && topSearch) {
      mobileSearchToggle.addEventListener('click', function () {
        topSearch.classList.toggle('mobile-open');
        const input = topSearch.querySelector('input');
        if (topSearch.classList.contains('mobile-open') && input) setTimeout(function () { input.focus(); }, 120);
      });
    }

    /* Focus search with ⌘K / Ctrl K or "/" */
    const searchInput = document.querySelector('.topbar-search input');
    const offcanvasEl = document.getElementById('searchOffcanvas');
    function openSearch() {
      if (offcanvasEl && window.bootstrap) {
        window.bootstrap.Offcanvas.getOrCreateInstance(offcanvasEl).show();
        setTimeout(function () {
          const input = offcanvasEl.querySelector('input');
          if (input) input.focus();
        }, 250);
      } else if (searchInput) {
        searchInput.focus();
      }
    }
    document.querySelectorAll('[data-open-search]').forEach(function (b) {
      b.addEventListener('click', openSearch);
    });
    document.addEventListener('keydown', function (e) {
      const tag = (e.target.tagName || '').toLowerCase();
      const typing = tag === 'input' || tag === 'textarea' || e.target.isContentEditable;
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') { e.preventDefault(); openSearch(); return; }
      if (!typing && e.key === '/') { e.preventDefault(); openSearch(); }
    });

    /* Search panel live filter (delegated, survives htmx swaps) */
    document.addEventListener('input', function (e) {
      if (!e.target.matches || !e.target.matches('#searchOffcanvasBody input[data-search-panel]')) return;
      const q = e.target.value.trim().toLowerCase();
      document.querySelectorAll('#searchOffcanvasBody [data-search-item]').forEach(function (a) {
        a.style.display = !q || a.innerText.toLowerCase().indexOf(q) !== -1 ? '' : 'none';
      });
    });

    /* Favorites (demo, stored locally) */
    const favs = (function () {
      try { return new Set(JSON.parse(localStorage.getItem('aitool-favs') || '[]')); } catch (e) { return new Set(); }
    })();
    function paintFavs() {
      document.querySelectorAll('[data-fav]').forEach(function (btn) {
        const on = favs.has(btn.dataset.fav);
        btn.classList.toggle('is-fav', on);
        btn.setAttribute('aria-pressed', on ? 'true' : 'false');
        const icon = btn.querySelector('use');
        if (icon) icon.setAttribute('href', '/assets/static/nav-001/icons.svg#' + (on ? 'heart' : 'heart'));
        btn.style.color = on ? 'var(--accent)' : '';
      });
    }
    document.addEventListener('click', function (e) {
      const btn = e.target.closest('[data-fav]');
      if (!btn) return;
      e.preventDefault();
      const key = btn.dataset.fav;
      if (favs.has(key)) favs.delete(key); else favs.add(key);
      try { localStorage.setItem('aitool-favs', JSON.stringify(Array.from(favs))); } catch (err) {}
      paintFavs();
    });
    paintFavs();

    /* Carousel */
    document.querySelectorAll('[data-carousel]').forEach(function (root) {
      const track = root.querySelector('[data-carousel-track]');
      if (!track) return;
      const step = function () { return Math.max(260, track.clientWidth * .35); };
      root.querySelectorAll('[data-carousel-prev]').forEach(function (b) { b.addEventListener('click', function () { track.scrollBy({ left: -step(), behavior: 'smooth' }); }); });
      root.querySelectorAll('[data-carousel-next]').forEach(function (b) { b.addEventListener('click', function () { track.scrollBy({ left: step(), behavior: 'smooth' }); }); });
    });

    /* Generic accordion / FAQ */
    document.addEventListener('click', function (e) {
      const q = e.target.closest('.faq-q');
      if (!q) return;
      const item = q.closest('.faq-item');
      if (!item) return;
      const isOpen = item.classList.contains('open');
      const parent = item.parentElement;
      if (parent) parent.querySelectorAll('.faq-item.open').forEach(function (x) { if (x !== item) x.classList.remove('open'); });
      item.classList.toggle('open', !isOpen);
      q.setAttribute('aria-expanded', String(!isOpen));
    });

    /* Filter tabs (tutorials, topics) */
    document.querySelectorAll('[data-filter-tab]').forEach(function (tab) {
      tab.addEventListener('click', function () {
        const group = tab.dataset.filterGroup;
        if (group) {
          document.querySelectorAll('[data-filter-group="' + group + '"] [data-filter-tab]').forEach(function (x) { x.classList.remove('active'); });
        }
        tab.classList.add('active');
        const kind = tab.dataset.filterTab;
        const target = tab.dataset.filterTarget || 'body';
        document.querySelectorAll(target + ' [data-kind]').forEach(function (el) {
          el.style.display = (kind === 'all' || el.dataset.kind === kind) ? '' : 'none';
        });
        document.querySelectorAll(target + ' [data-kind-section]').forEach(function (el) {
          const hasVisible = Array.from(el.querySelectorAll('[data-kind]')).some(function (x) { return x.style.display !== 'none'; });
          el.style.display = hasVisible ? '' : 'none';
        });
      });
    });

    /* Category filters (also works after htmx swap) */
    const categoryFilterState = { filter: 'all', sort: 'default' };
    function applyToolFilters() {
      const grid = document.getElementById('categoryToolGrid');
      if (!grid) return;
      const cards = Array.from(grid.querySelectorAll('[data-tool-card]'));
      cards.forEach(function (card) {
        const f = categoryFilterState.filter;
        let show = true;
        if (f === 'free') show = card.dataset.free === '1';
        if (f === 'tutorial') show = card.dataset.tutorial === '1';
        if (f === 'rated') show = parseFloat(card.dataset.rating || '0') >= 4;
        card.style.display = show ? '' : 'none';
      });
      if (categoryFilterState.sort === 'rating') {
        cards.slice().sort(function (a, b) { return parseFloat(b.dataset.rating || 0) - parseFloat(a.dataset.rating || 0); }).forEach(function (c) { grid.appendChild(c); });
      } else if (categoryFilterState.sort === 'name') {
        cards.slice().sort(function (a, b) { return (a.dataset.name || '').localeCompare(b.dataset.name || '', 'zh-CN'); }).forEach(function (c) { grid.appendChild(c); });
      } else if (categoryFilterState.sort === 'new') {
        cards.slice().sort(function (a, b) { return (b.dataset.fresh || '0') - (a.dataset.fresh || '0'); }).forEach(function (c) { grid.appendChild(c); });
      }
      const visible = cards.filter(function (c) { return c.style.display !== 'none'; }).length;
      const countEl = document.getElementById('categoryVisibleCount');
      if (countEl) countEl.textContent = visible;
    }
    function initCategoryTools() {
      const grid = document.getElementById('categoryToolGrid');
      if (!grid) return;
      document.querySelectorAll('[data-category-filter]').forEach(function (b) {
        b.classList.toggle('active', b.dataset.categoryFilter === categoryFilterState.filter);
      });
      document.querySelectorAll('[data-category-sort]').forEach(function (b) {
        b.classList.toggle('active', b.dataset.categorySort === categoryFilterState.sort);
      });
      applyToolFilters();
    }
    document.addEventListener('click', function (e) {
      const f = e.target.closest('[data-category-filter]');
      if (f) { categoryFilterState.filter = f.dataset.categoryFilter; initCategoryTools(); }
      const s = e.target.closest('[data-category-sort]');
      if (s) { categoryFilterState.sort = s.dataset.categorySort; initCategoryTools(); }
    });
    document.body.addEventListener('htmx:afterSwap', function (e) {
      if (e.target && (e.target.id === 'categoryToolGrid' || e.target.closest && e.target.closest('#categoryToolGrid'))) initCategoryTools();
    });

    /* Apps filters */
    const appQuery = (new URLSearchParams(location.search).get('q') || '').trim().toLowerCase();
    function applyAppFilter(which) {
      document.querySelectorAll('#appCategories .app-category').forEach(function (section) {
        let count = 0;
        section.querySelectorAll('.app-card').forEach(function (card) {
          let show = true;
          if (which === 'free') show = card.dataset.free === '1';
          if (which === 'tutorial') show = card.dataset.tutorial === '1';
          if (which === 'rated') show = parseFloat(card.dataset.rating || '0') >= 4;
          if (appQuery && card.innerText.toLowerCase().indexOf(appQuery) === -1) show = false;
          card.style.display = show ? '' : 'none';
          if (show) count++;
        });
        const counter = section.querySelector('[data-app-count]');
        if (counter) counter.textContent = count;
        section.style.display = count ? '' : 'none';
      });
    }
    document.querySelectorAll('[data-app-filter]').forEach(function (b) {
      b.addEventListener('click', function () {
        document.querySelectorAll('[data-app-filter]').forEach(function (x) { x.classList.remove('active'); });
        b.classList.add('active');
        applyAppFilter(b.dataset.appFilter);
      });
    });
    if (appQuery) applyAppFilter('all');
    document.addEventListener('click', function (e) {
      if (!e.target.closest || !e.target.closest('#loadMoreApps')) return;
      let tries = 0;
      const timer = setInterval(function () {
        tries++;
        const loaded = document.querySelectorAll('#appCategories .app-category').length >= 18;
        if (loaded || tries > 20) {
          const wrap = document.getElementById('loadMoreAppsWrap');
          if (wrap) wrap.remove();
          const btn = document.getElementById('loadMoreApps');
          if (btn) btn.remove();
          clearInterval(timer);
          if (loaded && location.hash.indexOf('#cat-') === 0) {
            const target = document.querySelector(location.hash);
            if (target) setTimeout(function () {
              target.scrollIntoView({ behavior: 'smooth', block: 'start' });
            }, 120);
          }
        }
      }, 150);
    });
    function scrollToAppHash() {
      if (!location.hash || !location.hash.startsWith('#cat-')) return;
      const target = document.querySelector(location.hash);
      if (target) { setTimeout(function () { target.scrollIntoView({ behavior: 'smooth', block: 'start' }); }, 80); }
      else {
        const btn = document.getElementById('loadMoreApps');
        if (btn) {
          setTimeout(function () { btn.click(); }, 300);
          setTimeout(function () {
            const el = document.querySelector(location.hash);
            if (el) el.scrollIntoView({ behavior: 'smooth', block: 'start' });
          }, 1400);
        }
      }
    }
    if (location.hash.startsWith('#cat-')) scrollToAppHash();
    window.addEventListener('hashchange', scrollToAppHash);

    /* Wiki alphabet + search */
    function initWiki() {
      const cards = Array.from(document.querySelectorAll('#wikiGrid [data-wiki-card]'));
      if (!cards.length) return;
      const q = (document.querySelector('#wikiSearchInput') || {}).value || '';
      const active = document.querySelector('[data-alpha].active');
      const letter = active ? active.dataset.alpha : 'all';
      cards.forEach(function (card) {
        const title = (card.dataset.title || '').toLowerCase();
        const matchLetter = letter === 'all' || (card.dataset.initial || '').toUpperCase() === letter;
        const matchSearch = !q || title.indexOf(q.toLowerCase()) !== -1 || (card.dataset.en || '').toLowerCase().indexOf(q.toLowerCase()) !== -1;
        card.style.display = (matchLetter && matchSearch) ? '' : 'none';
      });
    }
    document.querySelectorAll('[data-alpha]').forEach(function (b) {
      b.addEventListener('click', function () {
        document.querySelectorAll('[data-alpha]').forEach(function (x) { x.classList.remove('active'); });
        b.classList.add('active');
        initWiki();
      });
    });
    const wikiSearch = document.getElementById('wikiSearchInput');
    if (wikiSearch) wikiSearch.addEventListener('input', initWiki);

    /* Projects category filter */
    document.querySelectorAll('[data-project-filter]').forEach(function (b) {
      b.addEventListener('click', function () {
        document.querySelectorAll('[data-project-filter]').forEach(function (x) { x.classList.remove('active'); });
        b.classList.add('active');
        const cat = b.dataset.projectFilter;
        document.querySelectorAll('#projectList [data-project-card]').forEach(function (card) {
          card.style.display = (cat === 'all' || card.dataset.category === cat) ? '' : 'none';
        });
      });
    });

    /* Article TOC active state and smooth scroll */
    const tocLinks = Array.from(document.querySelectorAll('.article-toc a[href^="#"]'));
    if (tocLinks.length && 'IntersectionObserver' in window) {
      const targets = tocLinks.map(function (a) { return document.querySelector(a.getAttribute('href')); }).filter(Boolean);
      const io = new IntersectionObserver(function (entries) {
        entries.forEach(function (entry) {
          if (!entry.isIntersecting) return;
          tocLinks.forEach(function (a) { a.classList.toggle('active', a.getAttribute('href') === '#' + entry.target.id); });
        });
      }, { rootMargin: '-90px 0px -70% 0px', threshold: 0 });
      targets.forEach(function (t) { io.observe(t); });
    }
    tocLinks.forEach(function (a) {
      a.addEventListener('click', function (e) {
        const el = document.querySelector(a.getAttribute('href'));
        if (el) { e.preventDefault(); el.scrollIntoView({ behavior: 'smooth', block: 'start' }); history.replaceState(null, '', a.getAttribute('href')); }
      });
    });

    /* Share */
    document.querySelectorAll('[data-share]').forEach(function (btn) {
      btn.addEventListener('click', function () {
        const data = { title: document.title, url: location.href };
        if (navigator.share) navigator.share(data).catch(function () {});
        else if (navigator.clipboard) navigator.clipboard.writeText(location.href).then(function () {
          const old = btn.innerHTML; btn.textContent = '已复制'; setTimeout(function () { btn.innerHTML = old; }, 1600);
        }).catch(function () {});
      });
    });

    /* Detail tabs */
    document.querySelectorAll('[data-detail-tab]').forEach(function (tab) {
      tab.addEventListener('click', function () {
        document.querySelectorAll('[data-detail-tab]').forEach(function (x) { x.classList.remove('active'); });
        tab.classList.add('active');
        const name = tab.dataset.detailTab;
        const detail = document.getElementById('toolTabDetail');
        const reviews = document.getElementById('toolTabReviews');
        if (detail) detail.hidden = name !== 'detail';
        if (reviews) reviews.hidden = name !== 'reviews';
      });
    });

    /* Ranking period tabs */
    document.querySelectorAll('[data-rank-tab]').forEach(function (btn) {
      btn.addEventListener('click', function () {
        document.querySelectorAll('[data-rank-tab]').forEach(function (x) { x.classList.remove('active'); });
        btn.classList.add('active');
        const period = document.getElementById('rankPeriodText');
        if (period) period.textContent = btn.dataset.rankTab === 'month' ? '月榜 · 2026年10月 更新' : '周榜 · 2026年10月19日 更新';
      });
    });

    /* Compare hot tool search */
    const compareSearch = document.querySelector('[data-compare-search]');
    if (compareSearch) {
      compareSearch.addEventListener('input', function () {
        const q = compareSearch.value.trim().toLowerCase();
        document.querySelectorAll('#compareHotGrid .hot-compare-card').forEach(function (card) {
          card.style.display = (!q || card.innerText.toLowerCase().indexOf(q) !== -1) ? '' : 'none';
        });
      });
    }

    /* Topic search */
    const topicSearch = document.querySelector('[data-topic-search]');
    if (topicSearch) {
      topicSearch.addEventListener('input', function () {
        const q = topicSearch.value.trim().toLowerCase();
        document.querySelectorAll('.topic-card').forEach(function (card) {
          card.style.display = (!q || card.innerText.toLowerCase().indexOf(q) !== -1) ? '' : 'none';
        });
        document.querySelectorAll('[data-kind-section]').forEach(function (section) {
          const visible = Array.from(section.querySelectorAll('.topic-card')).some(function (c) { return c.style.display !== 'none'; });
          section.style.display = visible ? '' : 'none';
        });
      });
    }

    /* Compare */
    initCompare();

    /* Htmx utilities */
    document.body.addEventListener('htmx:afterSwap', function (e) {
      if (e.target && e.target.id === 'appCategories') {
        const active = document.querySelector('[data-app-filter].active');
        if (active) applyAppFilter(active.dataset.appFilter);
        const more = document.getElementById('loadMoreApps');
        if (more) more.remove();
        if (location.hash && location.hash.startsWith('#cat-')) setTimeout(scrollToAppHash, 50);
      }
      paintFavs();
    });

    // remove load more wrapper after successful request
    document.body.addEventListener('htmx:afterRequest', function (e) {
      if (e.target && e.target.id === 'loadMoreApps' && e.detail && e.detail.successful) {
        const wrap = document.getElementById('loadMoreAppsWrap');
        if (wrap) wrap.remove();
      }
    });
  }

  /* ---- Compare page ---- */
  function initCompare() {
    const page = document.getElementById('comparePage');
    if (!page) return;
    const tools = window.COMPARE_TOOLS || {};
    const chosen = [];
    const tray = document.getElementById('compareTray');
    const countEl = document.getElementById('compareCount');
    const statusEl = document.getElementById('compareStatusText');
    const result = document.getElementById('compareResult');

    function renderTray() {
      if (!tray) return;
      if (!chosen.length) {
        tray.innerHTML = '<div class="compare-tray-empty"><svg class="icon"><use href="/assets/static/nav-001/icons.svg#scale"/></svg>还没有选择工具，从下方热门工具里挑 2–4 款开始对比</div>';
      } else {
        tray.innerHTML = chosen.map(function (id) {
          const t = tools[id];
          return '<div class="compare-selected"><img src="' + t.icon + '" alt="' + t.name + '"><span>' + t.name + '</span><button type="button" data-remove="' + id + '" aria-label="移除 ' + t.name + '"><svg class="icon icon-sm"><use href="/assets/static/nav-001/icons.svg#x"/></svg></button></div>';
        }).join('');
      }
      if (countEl) countEl.textContent = chosen.length;
      if (statusEl) statusEl.textContent = chosen.length === 0 ? '至少再选 2 个就能开始对比' : (chosen.length < 2 ? '还差 ' + (2 - chosen.length) + ' 个就能开始对比' : (chosen.length === 4 ? '已选满 4 个，可以开始逐项对比' : '已选 ' + chosen.length + ' 个，可继续添加或开始对比'));
      document.querySelectorAll('[data-add-compare]').forEach(function (b) {
        const id = b.dataset.addCompare;
        b.classList.toggle('added', chosen.indexOf(id) !== -1);
        b.textContent = chosen.indexOf(id) !== -1 ? '已加入' : '加入对比';
      });
      if (!result) return;
      if (chosen.length < 2) {
        result.innerHTML = '';
        return;
      }
      const cols = chosen.map(function (id) { return tools[id]; });
      const rows = [
        ['分类', function (t) { return t.category; }],
        ['价格', function (t) { return t.price; }],
        ['评分', function (t) { return t.rating ? '★ ' + t.rating : '暂无评分'; }],
        ['核心能力', function (t) { return t.features; }],
        ['适合人群', function (t) { return t.audience; }],
        ['优点', function (t) { return t.pros; }],
        ['不足', function (t) { return t.cons; }]
      ];
      let html = '<div class="compare-table-wrap"><h2>逐项对比 <span class="text-secondary-soft" style="font-size:12px;font-weight:400">· 共 ' + cols.length + ' 款</span></h2><div class="table-responsive"><table class="compare-table"><thead><tr><th>对比项</th>';
      html += cols.map(function (t) {
        return '<th><div class="compare-tool-card"><img src="' + t.icon + '" alt="' + t.name + '"><div><h3>' + t.name + '</h3><p>' + t.desc + '</p></div></div></th>';
      }).join('');
      html += '</tr></thead><tbody>';
      rows.forEach(function (row) {
        html += '<tr><td>' + row[0] + '</td>';
        cols.forEach(function (t) { html += '<td>' + row[1](t) + '</td>'; });
        html += '</tr>';
      });
      html += '<tr><td>官网</td>' + cols.map(function (t) { return '<td><a class="tag" href="' + t.url + '" target="_blank" rel="noopener">访问官网 ↗</a></td>'; }).join('') + '</tr>';
      html += '</tbody></table></div></div>';
      result.innerHTML = html;
    }

    document.querySelectorAll('[data-add-compare]').forEach(function (btn) {
      btn.addEventListener('click', function (e) {
        e.preventDefault();
        const id = btn.dataset.addCompare;
        const idx = chosen.indexOf(id);
        if (idx !== -1) chosen.splice(idx, 1);
        else {
          if (chosen.length >= 4) { statusEl.textContent = '最多同时对比 4 款工具，请先移除一个'; return; }
          chosen.push(id);
        }
        renderTray();
      });
    });
    if (tray) tray.addEventListener('click', function (e) {
      const b = e.target.closest('[data-remove]');
      if (!b) return;
      const idx = chosen.indexOf(b.dataset.remove);
      if (idx !== -1) chosen.splice(idx, 1);
      renderTray();
    });
    const params = new URLSearchParams(location.search);
    if (params.get('tools')) {
      params.get('tools').split(',').slice(0, 4).forEach(function (id) { if (tools[id]) chosen.push(id); });
    }
    renderTray();
  }

  if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', init);
  else init();
})();
