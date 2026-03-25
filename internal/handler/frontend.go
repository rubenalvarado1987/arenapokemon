package handler

import (
	"net/http"
	"strings"
)

const frontendHTML = `<!DOCTYPE html>
<html lang="es">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>__APP_NAME__ Pokedex</title>
  <style>
    :root {
      --bg: #f4efe6;
      --surface: rgba(255, 251, 243, 0.82);
      --surface-strong: #fff7eb;
      --surface-dark: rgba(33, 29, 28, 0.9);
      --text: #241d1c;
      --muted: #6f655e;
      --brand: #ef476f;
      --brand-dark: #be3455;
      --accent: #ffd166;
      --accent-2: #06d6a0;
      --line: rgba(36, 29, 28, 0.08);
      --shadow: 0 24px 70px rgba(70, 45, 20, 0.14);
      --radius: 26px;
    }

    * { box-sizing: border-box; }

    body {
      margin: 0;
      min-height: 100vh;
      font-family: Georgia, "Times New Roman", serif;
      color: var(--text);
      background:
        radial-gradient(circle at top left, rgba(255, 209, 102, 0.74), transparent 24%),
        radial-gradient(circle at 85%% 10%%, rgba(239, 71, 111, 0.18), transparent 22%),
        linear-gradient(180deg, #fff7ea 0%%, var(--bg) 56%%, #ebe0cf 100%%);
    }

    .page {
      width: min(1200px, calc(100%% - 32px));
      margin: 0 auto;
      padding: 40px 0 56px;
    }

    .hero {
      display: grid;
      gap: 16px;
      padding: 32px;
      border: 1px solid var(--line);
      border-radius: calc(var(--radius) + 8px);
      background: linear-gradient(135deg, rgba(255, 249, 237, 0.96), rgba(255, 255, 255, 0.75));
      box-shadow: var(--shadow);
      overflow: hidden;
      position: relative;
    }

    .hero::after {
      content: "";
      position: absolute;
      inset: auto -40px -50px auto;
      width: 220px;
      height: 220px;
      border-radius: 50%%;
      background: radial-gradient(circle, rgba(239, 71, 111, 0.18), rgba(239, 71, 111, 0));
    }

    h1 {
      margin: 0;
      font-size: clamp(2.4rem, 5vw, 4.8rem);
      line-height: 0.92;
      letter-spacing: -0.04em;
      max-width: 9ch;
    }

    .hero p {
      margin: 0;
      max-width: 62ch;
      color: var(--muted);
      font-size: 1rem;
    }

    .hero-top {
      display: flex;
      justify-content: space-between;
      gap: 18px;
      align-items: start;
      flex-wrap: wrap;
    }

    .hero-stats {
      display: grid;
      grid-template-columns: repeat(3, minmax(96px, 1fr));
      gap: 10px;
      min-width: min(100%%, 320px);
    }

    .stat-chip {
      padding: 12px 14px;
      border-radius: 18px;
      background: rgba(255, 255, 255, 0.74);
      border: 1px solid rgba(36, 29, 28, 0.08);
    }

    .stat-chip strong,
    .stat-chip span {
      display: block;
    }

    .stat-chip span {
      color: var(--muted);
      font-size: 0.8rem;
      margin-bottom: 4px;
    }

    .toolbar {
      display: flex;
      gap: 12px;
      flex-wrap: wrap;
      align-items: center;
      margin-top: 8px;
    }

    .toolbar input {
      flex: 1 1 220px;
      min-width: 0;
      border: 1px solid rgba(32, 32, 32, 0.12);
      border-radius: 999px;
      padding: 14px 18px;
      background: rgba(255, 255, 255, 0.9);
      font: inherit;
    }

    .toolbar button {
      border: 0;
      border-radius: 999px;
      padding: 14px 18px;
      background: var(--brand);
      color: #fff;
      font: inherit;
      cursor: pointer;
      transition: transform 160ms ease, background 160ms ease;
    }

    .toolbar button:hover {
      background: var(--brand-dark);
      transform: translateY(-1px);
    }

    .arena {
      margin: 26px 0 22px;
      display: grid;
      gap: 16px;
      grid-template-columns: 1.2fr 0.8fr;
      align-items: stretch;
    }

    .battle-stage,
    .battle-log {
      border-radius: calc(var(--radius) + 4px);
      border: 1px solid var(--line);
      box-shadow: var(--shadow);
      overflow: hidden;
    }

    .battle-stage {
      background:
        linear-gradient(180deg, rgba(255, 245, 224, 0.9), rgba(255, 255, 255, 0.72)),
        linear-gradient(120deg, rgba(239, 71, 111, 0.08), rgba(6, 214, 160, 0.08));
      padding: 24px;
      display: grid;
      gap: 20px;
    }

    .battle-header {
      display: flex;
      justify-content: space-between;
      gap: 12px;
      align-items: center;
      flex-wrap: wrap;
    }

    .battle-header h2 {
      margin: 0;
      font-size: 1.6rem;
    }

    .battle-subtitle {
      margin: 4px 0 0;
      color: var(--muted);
      font-size: 0.96rem;
    }

    .battle-actions {
      display: flex;
      gap: 10px;
      flex-wrap: wrap;
    }

    .battle-actions button {
      border: 0;
      border-radius: 999px;
      padding: 12px 16px;
      font: inherit;
      cursor: pointer;
    }

    .battle-actions .primary {
      background: var(--surface-dark);
      color: #fff;
    }

    .battle-actions .secondary {
      background: rgba(255, 255, 255, 0.76);
      color: var(--text);
      border: 1px solid rgba(36, 29, 28, 0.08);
    }

    .fighters {
      display: grid;
      grid-template-columns: 1fr auto 1fr;
      gap: 16px;
      align-items: center;
    }

    .fighter {
      padding: 18px;
      border-radius: 24px;
      background: rgba(255, 255, 255, 0.62);
      border: 1px solid rgba(36, 29, 28, 0.08);
      display: grid;
      gap: 12px;
      position: relative;
      min-height: 300px;
    }

    .fighter.active {
      animation: pulse 900ms ease-in-out infinite alternate;
      box-shadow: 0 0 0 1px rgba(239, 71, 111, 0.22), 0 18px 40px rgba(239, 71, 111, 0.12);
    }

    .fighter.winner {
      background: linear-gradient(180deg, rgba(6, 214, 160, 0.18), rgba(255, 255, 255, 0.72));
    }

    .fighter.loser {
      opacity: 0.55;
      filter: grayscale(0.25);
    }

    .fighter-avatar {
      display: grid;
      place-items: center;
      min-height: 160px;
      border-radius: 20px;
      background: radial-gradient(circle at top, rgba(255, 209, 102, 0.28), rgba(255, 255, 255, 0.6));
      position: relative;
      overflow: hidden;
    }

    .fighter-avatar::after {
      content: "";
      position: absolute;
      inset: auto 12%% 8%% 12%%;
      height: 18px;
      background: rgba(36, 29, 28, 0.08);
      filter: blur(12px);
      border-radius: 50%%;
    }

    .fighter-avatar img {
      width: min(180px, 75%%);
      aspect-ratio: 1;
      object-fit: contain;
      z-index: 1;
      transition: transform 220ms ease, filter 220ms ease;
    }

    .fighter.active .fighter-avatar img {
      transform: translateY(-6px) scale(1.04);
    }

    .fighter.hit .fighter-avatar img {
      animation: shake 320ms ease;
      filter: saturate(1.2);
    }

    .fighter-header {
      display: flex;
      justify-content: space-between;
      gap: 8px;
      align-items: start;
    }

    .fighter-header strong {
      display: block;
      font-size: 1.4rem;
      text-transform: capitalize;
    }

    .fighter-header span,
    .fighter-meta {
      color: var(--muted);
      font-size: 0.9rem;
    }

    .fighter-stats {
      display: grid;
      gap: 8px;
    }

    .meter {
      display: grid;
      gap: 6px;
    }

    .meter-head {
      display: flex;
      justify-content: space-between;
      gap: 8px;
      font-size: 0.86rem;
      color: var(--muted);
    }

    .meter-track {
      height: 12px;
      border-radius: 999px;
      background: rgba(36, 29, 28, 0.08);
      overflow: hidden;
    }

    .meter-fill {
      height: 100%%;
      width: 100%%;
      border-radius: 999px;
      background: linear-gradient(90deg, var(--brand), var(--accent));
      transition: width 340ms ease;
    }

    .versus {
      display: grid;
      place-items: center;
      min-width: 72px;
      color: #fff;
      font-weight: 700;
      letter-spacing: 0.12em;
      text-transform: uppercase;
    }

    .versus-badge {
      width: 70px;
      height: 70px;
      border-radius: 50%%;
      background: radial-gradient(circle at 30%% 30%%, #ff7b9c, var(--brand-dark));
      display: grid;
      place-items: center;
      box-shadow: 0 18px 30px rgba(190, 52, 85, 0.3);
    }

    .battle-event {
      border-radius: 22px;
      padding: 18px 20px;
      background: rgba(255, 255, 255, 0.82);
      border: 1px solid rgba(36, 29, 28, 0.08);
      display: grid;
      gap: 8px;
    }

    .battle-event strong {
      font-size: 1.12rem;
    }

    .battle-event p,
    .battle-log p {
      margin: 0;
      color: var(--muted);
    }

    .battle-log {
      background: linear-gradient(180deg, rgba(32, 29, 28, 0.95), rgba(48, 40, 38, 0.96));
      color: #fff8ee;
      padding: 22px;
      display: grid;
      gap: 14px;
      align-content: start;
    }

    .battle-log h3,
    .rounds h3,
    .gallery-header h3 {
      margin: 0;
      font-size: 1.1rem;
    }

    .battle-log-list,
    .round-list {
      display: grid;
      gap: 10px;
      max-height: 360px;
      overflow: auto;
      padding-right: 4px;
    }

    .battle-log-item,
    .round-item {
      padding: 12px 14px;
      border-radius: 16px;
      background: rgba(255, 255, 255, 0.07);
      border: 1px solid rgba(255, 255, 255, 0.07);
    }

    .battle-log-item small,
    .round-item small {
      display: block;
      color: rgba(255, 248, 238, 0.68);
      margin-top: 4px;
    }

    .rounds {
      margin: 0 0 26px;
      display: grid;
      gap: 12px;
      padding: 22px;
      border-radius: calc(var(--radius) + 4px);
      background: rgba(255, 255, 255, 0.66);
      border: 1px solid var(--line);
      box-shadow: var(--shadow);
    }

    .gallery-header {
      display: flex;
      justify-content: space-between;
      gap: 12px;
      align-items: center;
      flex-wrap: wrap;
    }

    .meta {
      display: flex;
      justify-content: space-between;
      gap: 12px;
      align-items: center;
      margin: 22px 0 16px;
      color: var(--muted);
      font-size: 0.96rem;
    }

    .grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
      gap: 18px;
    }

    .card {
      background: var(--surface);
      border: 1px solid var(--line);
      border-radius: var(--radius);
      padding: 18px;
      box-shadow: var(--shadow);
      backdrop-filter: blur(10px);
      display: grid;
      gap: 12px;
      transform: translateY(18px);
      opacity: 0;
      animation: rise 500ms ease forwards;
    }

    .card img {
      width: 100%%;
      aspect-ratio: 1;
      object-fit: contain;
      border-radius: 18px;
      background: linear-gradient(180deg, rgba(255, 209, 102, 0.18), rgba(255, 255, 255, 0.7));
      padding: 12px;
    }

    .card-title {
      display: flex;
      justify-content: space-between;
      gap: 8px;
      align-items: baseline;
    }

    .card-title strong {
      text-transform: capitalize;
      font-size: 1.1rem;
    }

    .badge {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      border-radius: 999px;
      padding: 6px 10px;
      background: var(--surface-strong);
      color: var(--brand-dark);
      font-size: 0.82rem;
      border: 1px solid rgba(239, 71, 111, 0.18);
    }

    .empty {
      padding: 28px;
      border-radius: var(--radius);
      background: rgba(255, 255, 255, 0.7);
      border: 1px dashed rgba(32, 32, 32, 0.16);
      text-align: center;
      color: var(--muted);
    }

    .sparkles {
      position: absolute;
      inset: 0;
      pointer-events: none;
      overflow: hidden;
    }

    .spark {
      position: absolute;
      width: 12px;
      height: 12px;
      border-radius: 50%%;
      background: radial-gradient(circle, rgba(255,255,255,0.94), rgba(255,209,102,0.2));
      animation: spark 900ms ease forwards;
    }

    @keyframes rise {
      from { opacity: 0; transform: translateY(18px); }
      to { opacity: 1; transform: translateY(0); }
    }

    @keyframes pulse {
      from { transform: translateY(0); }
      to { transform: translateY(-4px); }
    }

    @keyframes shake {
      0%% { transform: translateX(0); }
      20%% { transform: translateX(-8px); }
      40%% { transform: translateX(7px); }
      60%% { transform: translateX(-5px); }
      80%% { transform: translateX(4px); }
      100%% { transform: translateX(0); }
    }

    @keyframes spark {
      from { transform: translateY(0) scale(0.6); opacity: 0.9; }
      to { transform: translateY(-60px) scale(1.6); opacity: 0; }
    }

    @media (max-width: 640px) {
      .page { width: min(100%% - 20px, 1200px); padding-top: 20px; }
      .hero { padding: 22px; }
      .meta { align-items: flex-start; flex-direction: column; }
    }

    @media (max-width: 960px) {
      .arena { grid-template-columns: 1fr; }
      .fighters { grid-template-columns: 1fr; }
      .versus { order: -1; }
    }
  </style>
</head>
<body>
  <main class="page">
    <section class="hero">
      <div class="hero-top">
        <div>
          <span class="badge">Campeonato automático</span>
          <h1>Pokedex arena en vivo</h1>
          <p>La página arranca con una batalla aleatoria, aplica efectos visuales, elige un ganador por rondas y continúa hasta coronar un campeón. Debajo se mantiene la galería completa en cards.</p>
        </div>
        <div class="hero-stats">
          <div class="stat-chip"><span>Formato</span><strong id="tournamentFormat">Cargando...</strong></div>
          <div class="stat-chip"><span>Batalla actual</span><strong id="battleNumber">--</strong></div>
          <div class="stat-chip"><span>Campeón</span><strong id="championName">Pendiente</strong></div>
        </div>
      </div>
      <div class="toolbar">
        <input id="search" type="search" placeholder="Filtrar por nombre" aria-label="Filtrar por nombre" />
        <button id="loadMore" type="button">Cargar más</button>
      </div>
    </section>

    <section class="arena">
      <article class="battle-stage">
        <div class="battle-header">
          <div>
            <h2 id="roundTitle">Preparando combate...</h2>
            <p class="battle-subtitle" id="battleSubtitle">Esperando suficientes Pokémon para iniciar el campeonato.</p>
          </div>
          <div class="battle-actions">
            <button class="primary" id="restartTournament" type="button">Nuevo campeonato</button>
            <button class="secondary" id="toggleAutoBattle" type="button">Pausar autoplay</button>
          </div>
        </div>

        <div class="fighters">
          <article class="fighter" id="fighterA">
            <div class="sparkles" id="sparklesA"></div>
            <div class="fighter-avatar"><img id="fighterAImage" alt="" /></div>
            <div class="fighter-header">
              <div>
                <strong id="fighterAName">---</strong>
                <span id="fighterASeed">Semilla --</span>
              </div>
              <span class="badge" id="fighterAState">Listo</span>
            </div>
            <div class="fighter-meta" id="fighterAMotto">Esperando rival.</div>
            <div class="fighter-stats">
              <div class="meter">
                <div class="meter-head"><span>Energía</span><span id="fighterAEnergyLabel">100 / 100</span></div>
                <div class="meter-track"><div class="meter-fill" id="fighterAEnergyBar"></div></div>
              </div>
              <div class="meter">
                <div class="meter-head"><span>Momentum</span><span id="fighterAMomentumLabel">0</span></div>
                <div class="meter-track"><div class="meter-fill" id="fighterAMomentumBar"></div></div>
              </div>
            </div>
          </article>

          <div class="versus">
            <div class="versus-badge">VS</div>
          </div>

          <article class="fighter" id="fighterB">
            <div class="sparkles" id="sparklesB"></div>
            <div class="fighter-avatar"><img id="fighterBImage" alt="" /></div>
            <div class="fighter-header">
              <div>
                <strong id="fighterBName">---</strong>
                <span id="fighterBSeed">Semilla --</span>
              </div>
              <span class="badge" id="fighterBState">Listo</span>
            </div>
            <div class="fighter-meta" id="fighterBMotto">Esperando rival.</div>
            <div class="fighter-stats">
              <div class="meter">
                <div class="meter-head"><span>Energía</span><span id="fighterBEnergyLabel">100 / 100</span></div>
                <div class="meter-track"><div class="meter-fill" id="fighterBEnergyBar"></div></div>
              </div>
              <div class="meter">
                <div class="meter-head"><span>Momentum</span><span id="fighterBMomentumLabel">0</span></div>
                <div class="meter-track"><div class="meter-fill" id="fighterBMomentumBar"></div></div>
              </div>
            </div>
          </article>
        </div>

        <section class="battle-event">
          <strong id="battleEventTitle">Los aspirantes se reúnen en la arena.</strong>
          <p id="battleEventBody">En cuanto haya suficientes Pokémon cargados, comenzará un cuadro eliminatorio automático.</p>
        </section>
      </article>

      <aside class="battle-log">
        <div>
          <h3>Crónica del combate</h3>
          <p id="logStatus">Se mostrarán golpes, críticos y remontadas en tiempo real.</p>
        </div>
        <div id="battleLog" class="battle-log-list"></div>
      </aside>
    </section>

    <section class="rounds">
      <div class="gallery-header">
        <h3>Bracket del campeonato</h3>
        <span class="badge" id="remainingCount">0 en carrera</span>
      </div>
      <div id="roundHistory" class="round-list"></div>
    </section>

    <section class="meta">
      <div id="status">Cargando Pokémon...</div>
      <div id="counter"></div>
    </section>

    <div class="gallery-header">
      <h3>Pokémon disponibles</h3>
      <span class="badge">Galería de competidores</span>
    </div>
    <section id="grid" class="grid" aria-live="polite"></section>
  </main>

  <script>
    const pageSize = 24;
    const tournamentSize = 8;
    const maxBattleLog = 8;
    let offset = 0;
    let total = Infinity;
    let allPokemon = [];
    let tournament = null;
    let autoBattle = true;
    let pendingBattleTimeout = null;

    const grid = document.getElementById('grid');
    const status = document.getElementById('status');
    const counter = document.getElementById('counter');
    const search = document.getElementById('search');
    const loadMore = document.getElementById('loadMore');
    const roundTitle = document.getElementById('roundTitle');
    const battleSubtitle = document.getElementById('battleSubtitle');
    const battleEventTitle = document.getElementById('battleEventTitle');
    const battleEventBody = document.getElementById('battleEventBody');
    const battleLog = document.getElementById('battleLog');
    const logStatus = document.getElementById('logStatus');
    const roundHistory = document.getElementById('roundHistory');
    const remainingCount = document.getElementById('remainingCount');
    const tournamentFormat = document.getElementById('tournamentFormat');
    const battleNumber = document.getElementById('battleNumber');
    const championName = document.getElementById('championName');
    const restartTournament = document.getElementById('restartTournament');
    const toggleAutoBattle = document.getElementById('toggleAutoBattle');

    const fighterAElements = getFighterElements('A');
    const fighterBElements = getFighterElements('B');

    function getFighterElements(slot) {
      return {
        root: document.getElementById('fighter' + slot),
        image: document.getElementById('fighter' + slot + 'Image'),
        name: document.getElementById('fighter' + slot + 'Name'),
        seed: document.getElementById('fighter' + slot + 'Seed'),
        state: document.getElementById('fighter' + slot + 'State'),
        motto: document.getElementById('fighter' + slot + 'Motto'),
        energyLabel: document.getElementById('fighter' + slot + 'EnergyLabel'),
        energyBar: document.getElementById('fighter' + slot + 'EnergyBar'),
        momentumLabel: document.getElementById('fighter' + slot + 'MomentumLabel'),
        momentumBar: document.getElementById('fighter' + slot + 'MomentumBar'),
        sparkles: document.getElementById('sparkles' + slot)
      };
    }

    function render(items) {
      if (!items.length) {
        grid.innerHTML = '<div class="empty">No hay Pokémon que coincidan con el filtro actual.</div>';
        return;
      }

      grid.innerHTML = items.map(function(pokemon, index) {
        return [
          '<article class="card" style="animation-delay: ' + (index * 24) + 'ms">',
          '  <img src="' + (pokemon.image_url || '') + '" alt="' + pokemon.name + '" loading="lazy" />',
          '  <div class="card-title">',
          '    <strong>' + pokemon.name + '</strong>',
          '    <span class="badge">#' + indexOfPokemon(pokemon.name) + '</span>',
          '  </div>',
          '</article>'
        ].join('');
      }).join('');
    }

    function indexOfPokemon(name) {
      const index = allPokemon.findIndex((pokemon) => pokemon.name === name);
      return index === -1 ? '?' : index + 1;
    }

    function updateMeta(filteredCount) {
      const loaded = allPokemon.length;
      status.textContent = loaded >= total ? 'Todos los Pokémon disponibles han sido cargados.' : 'Explorando la Pokédex desde la API...';
      counter.textContent = filteredCount + ' visibles / ' + loaded + ' cargados' + (Number.isFinite(total) ? ' / ' + total + ' total' : '');
      loadMore.disabled = loaded >= total;
      loadMore.textContent = loadMore.disabled ? 'Sin más resultados' : 'Cargar más';
      ensureTournamentReady();
    }

    async function fetchPage() {
      loadMore.disabled = true;
      status.textContent = 'Cargando Pokémon...';

      const response = await fetch('/api/v1/pokemon?limit=' + pageSize + '&offset=' + offset);
      if (!response.ok) {
        throw new Error('No se pudo cargar la lista de Pokémon');
      }

      const payload = await response.json();
      total = payload.count;
      offset += payload.results.length;
      allPokemon = allPokemon.concat(payload.results);
      applyFilter();
    }

    function applyFilter() {
      const term = search.value.trim().toLowerCase();
      const filtered = term ? allPokemon.filter((pokemon) => pokemon.name.includes(term)) : allPokemon;
      render(filtered);
      updateMeta(filtered.length);
    }

    function ensureTournamentReady() {
      if (tournament || allPokemon.length < tournamentSize) {
        tournamentFormat.textContent = allPokemon.length < tournamentSize ? 'Esperando ' + (tournamentSize - allPokemon.length) : tournamentSize + ' finalistas';
        return;
      }

      startTournament();
    }

    function startTournament() {
      clearPendingBattle();
      const contestants = shuffle(allPokemon.slice()).slice(0, tournamentSize).map(function(pokemon, index) {
        return {
          name: pokemon.name,
          image_url: pokemon.image_url,
          seed: index + 1,
          power: 72 + Math.floor(Math.random() * 29),
          flair: randomOf([
            'entra con aura electrica',
            'responde con una defensa de hierro',
            'invoca un sprint fulminante',
            'presiona con una tecnica impredecible',
            'mantiene una calma tactica absoluta'
          ])
        };
      });

      tournament = {
        contestants: contestants,
        queue: contestants.slice(),
        nextRoundQueue: [],
        round: 1,
        battle: 1,
        champion: null,
        history: []
      };

      tournamentFormat.textContent = contestants.length + ' competidores';
      championName.textContent = 'Pendiente';
      roundTitle.textContent = 'Cuartos de final listos';
      battleSubtitle.textContent = 'El campeonato comienza apenas se presenta el primer duelo.';
      battleEventTitle.textContent = 'Se definieron los participantes del campeonato.';
      battleEventBody.textContent = 'Los cruces son aleatorios y cada combate se resuelve con energia, momentum y un pequeño factor sorpresa.';
      battleLog.innerHTML = '';
      logStatus.textContent = 'Autoplay ' + (autoBattle ? 'activo' : 'pausado') + '. Las batallas se ejecutan en secuencia.';
      renderRoundHistory();
      runNextBattle();
    }

    function runNextBattle() {
      if (!tournament) {
        return;
      }

      if (tournament.queue.length === 1 && tournament.nextRoundQueue.length === 0) {
        announceChampion(tournament.queue[0]);
        return;
      }

      if (tournament.queue.length < 2) {
        tournament.queue = tournament.nextRoundQueue.slice();
        tournament.nextRoundQueue = [];
        tournament.round += 1;
        tournament.battle = 1;
        renderRoundHistory();
      }

      if (tournament.queue.length < 2) {
        announceChampion(tournament.queue[0]);
        return;
      }

      const left = tournament.queue.shift();
      const right = tournament.queue.shift();
      simulateBattle(left, right);
    }

    function simulateBattle(leftBase, rightBase) {
      const left = createBattleState(leftBase);
      const right = createBattleState(rightBase);
      const phase = roundLabel(tournament.queue.length + tournament.nextRoundQueue.length + 2);

      battleNumber.textContent = String(tournament.battle);
      roundTitle.textContent = phase;
      battleSubtitle.textContent = left.name + ' y ' + right.name + ' entran a la arena.';
      battleEventTitle.textContent = 'Comienza ' + left.name + ' vs ' + right.name + '.';
      battleEventBody.textContent = 'La pelea se resuelve en varios intercambios; el que vacie primero la energia rival avanza.';

      renderFighter(fighterAElements, left, 'En guardia');
      renderFighter(fighterBElements, right, 'En guardia');

      const turnLog = [];
      const totalTurns = 4 + Math.floor(Math.random() * 3);
      let attacker = left;
      let defender = right;

      for (let turn = 1; turn <= totalTurns; turn += 1) {
        const damage = calculateDamage(attacker, defender, turn);
        defender.energy = Math.max(0, defender.energy - damage);
        attacker.momentum = Math.min(100, attacker.momentum + 18 + Math.floor(Math.random() * 8));
        defender.momentum = Math.max(0, defender.momentum - 8);

        turnLog.push({
          attacker: attacker.name,
          defender: defender.name,
          damage: damage,
          remaining: defender.energy,
          text: attacker.name + ' golpea con ' + attacker.flair + ' y deja a ' + defender.name + ' en ' + defender.energy + ' de energia.'
        });

        if (defender.energy === 0) {
          break;
        }

        const tmp = attacker;
        attacker = defender;
        defender = tmp;
      }

      const winner = left.energy === right.energy ? tieBreaker(left, right) : (left.energy > right.energy ? left : right);
      const loser = winner.name === left.name ? right : left;
      tournament.nextRoundQueue.push(stripBattleState(winner));
      tournament.history.unshift({
        round: phase,
        battle: tournament.battle,
        winner: winner.name,
        loser: loser.name,
        summary: winner.name + ' elimina a ' + loser.name + ' con ' + winner.energy + ' de energia restante.'
      });
      tournament.history = tournament.history.slice(0, 12);

      playBattleAnimation(left, right, winner, loser, turnLog);
      tournament.battle += 1;
    }

    function playBattleAnimation(left, right, winner, loser, turnLog) {
      let step = 0;

      function nextStep() {
        if (step < turnLog.length) {
          const event = turnLog[step];
          const attackerElements = event.attacker === left.name ? fighterAElements : fighterBElements;
          const defenderElements = event.defender === left.name ? fighterAElements : fighterBElements;
          const attackerState = event.attacker === left.name ? left : right;
          const defenderState = event.defender === left.name ? left : right;

          attackerElements.root.classList.add('active');
          defenderElements.root.classList.add('hit');
          addSpark(defenderElements.sparkles);
          renderFighter(attackerElements, attackerState, 'Atacando');
          renderFighter(defenderElements, defenderState, 'Resistiendo');
          battleEventTitle.textContent = event.attacker + ' conecta un golpe.';
          battleEventBody.textContent = event.text;
          pushBattleLog(event.attacker + ' hace ' + event.damage + ' de daño a ' + event.defender + '.', event.remaining + ' energia restante');

          window.setTimeout(function() {
            attackerElements.root.classList.remove('active');
            defenderElements.root.classList.remove('hit');
          }, 260);

          step += 1;
          pendingBattleTimeout = window.setTimeout(nextStep, 760);
          return;
        }

        finalizeBattle(left, right, winner, loser);
      }

      nextStep();
    }

    function finalizeBattle(left, right, winner, loser) {
      const winnerElements = winner.name === left.name ? fighterAElements : fighterBElements;
      const loserElements = loser.name === left.name ? fighterAElements : fighterBElements;

      winnerElements.root.classList.add('winner');
      loserElements.root.classList.add('loser');
      renderFighter(winnerElements, winner, 'Ganador');
      renderFighter(loserElements, loser, 'Eliminado');
      battleEventTitle.textContent = winner.name + ' avanza a la siguiente ronda.';
      battleEventBody.textContent = winner.name + ' vence a ' + loser.name + ' y mantiene vivo su camino al titulo.';
      pushBattleLog(winner.name + ' gana la batalla ' + (tournament.battle), loser.name + ' queda fuera del cuadro');
      renderRoundHistory();

      if (!autoBattle) {
        logStatus.textContent = 'Autoplay pausado. Reanuda para continuar el campeonato.';
        return;
      }

      pendingBattleTimeout = window.setTimeout(function() {
        resetFighterClasses();
        runNextBattle();
      }, 1800);
    }

    function announceChampion(champion) {
      tournament.champion = champion;
      championName.textContent = champion.name;
      roundTitle.textContent = 'Campeon del torneo';
      battleSubtitle.textContent = champion.name + ' supera todo el bracket y se queda con la copa.';
      battleEventTitle.textContent = champion.name + ' es el nuevo campeon.';
      battleEventBody.textContent = 'El campeonato termino. Puedes iniciar otro cuadro para ver nuevos cruces y desenlaces.';
      pushBattleLog(champion.name + ' se corona campeon.', 'El cuadro completo se reiniciara cuando quieras.');
      renderRoundHistory();
      logStatus.textContent = 'Campeonato completado.';
    }

    function renderFighter(elements, fighter, stateLabel) {
      elements.name.textContent = fighter.name;
      elements.seed.textContent = 'Semilla ' + fighter.seed;
      elements.state.textContent = stateLabel;
      elements.motto.textContent = fighter.flair;
      elements.image.src = fighter.image_url || '';
      elements.image.alt = fighter.name;
      elements.energyLabel.textContent = fighter.energy + ' / 100';
      elements.energyBar.style.width = fighter.energy + '%';
      elements.momentumLabel.textContent = String(fighter.momentum);
      elements.momentumBar.style.width = fighter.momentum + '%';
    }

    function resetFighterClasses() {
      [fighterAElements.root, fighterBElements.root].forEach(function(node) {
        node.classList.remove('active');
        node.classList.remove('hit');
        node.classList.remove('winner');
        node.classList.remove('loser');
      });
    }

    function pushBattleLog(title, subtitle) {
      const item = document.createElement('div');
      item.className = 'battle-log-item';
      item.innerHTML = '<strong>' + title + '</strong><small>' + subtitle + '</small>';
      battleLog.prepend(item);
      while (battleLog.children.length > maxBattleLog) {
        battleLog.removeChild(battleLog.lastChild);
      }
    }

    function renderRoundHistory() {
      if (!tournament) {
        roundHistory.innerHTML = '<div class="round-item">Esperando participantes suficientes para el cuadro.</div>';
        remainingCount.textContent = '0 en carrera';
        return;
      }

      const alive = tournament.queue.length + tournament.nextRoundQueue.length;
      remainingCount.textContent = alive + ' en carrera';

      if (!tournament.history.length) {
        roundHistory.innerHTML = '<div class="round-item">Aun no se han resuelto combates.</div>';
        return;
      }

      roundHistory.innerHTML = tournament.history.map(function(entry) {
        return '<div class="round-item"><strong>' + entry.round + ' · Batalla ' + entry.battle + '</strong><small>' + entry.summary + '</small></div>';
      }).join('');
    }

    function createBattleState(base) {
      return {
        name: base.name,
        image_url: base.image_url,
        seed: base.seed,
        flair: base.flair,
        power: base.power,
        energy: 100,
        momentum: 16 + Math.floor(Math.random() * 18)
      };
    }

    function stripBattleState(fighter) {
      return {
        name: fighter.name,
        image_url: fighter.image_url,
        seed: fighter.seed,
        power: Math.min(100, fighter.power + 2),
        flair: fighter.flair
      };
    }

    function calculateDamage(attacker, defender, turn) {
      const swing = Math.floor(Math.random() * 12);
      const critical = Math.random() > 0.82 ? 14 : 0;
      const turnPressure = turn * 2;
      return Math.max(12, Math.min(44, Math.floor(attacker.power * 0.16) + Math.floor(attacker.momentum * 0.12) - Math.floor(defender.power * 0.05) + swing + critical + turnPressure));
    }

    function tieBreaker(left, right) {
      return left.momentum >= right.momentum ? left : right;
    }

    function roundLabel(participants) {
      if (participants <= 2) {
        return 'Gran final';
      }
      if (participants <= 4) {
        return 'Semifinales';
      }
      if (participants <= 8) {
        return 'Cuartos de final';
      }
      return 'Ronda ' + tournament.round;
    }

    function randomOf(list) {
      return list[Math.floor(Math.random() * list.length)];
    }

    function shuffle(list) {
      for (let i = list.length - 1; i > 0; i -= 1) {
        const j = Math.floor(Math.random() * (i + 1));
        const tmp = list[i];
        list[i] = list[j];
        list[j] = tmp;
      }
      return list;
    }

    function addSpark(container) {
      const spark = document.createElement('span');
      spark.className = 'spark';
      spark.style.left = (16 + Math.random() * 68) + '%';
      spark.style.top = (28 + Math.random() * 32) + '%';
      container.appendChild(spark);
      window.setTimeout(function() {
        if (spark.parentNode) {
          spark.parentNode.removeChild(spark);
        }
      }, 920);
    }

    function clearPendingBattle() {
      if (pendingBattleTimeout) {
        window.clearTimeout(pendingBattleTimeout);
        pendingBattleTimeout = null;
      }
      resetFighterClasses();
    }

    search.addEventListener('input', applyFilter);
    loadMore.addEventListener('click', () => {
      fetchPage().catch((error) => {
        status.textContent = error.message;
      });
    });
    restartTournament.addEventListener('click', startTournament);
    toggleAutoBattle.addEventListener('click', function() {
      autoBattle = !autoBattle;
      toggleAutoBattle.textContent = autoBattle ? 'Pausar autoplay' : 'Reanudar autoplay';
      logStatus.textContent = 'Autoplay ' + (autoBattle ? 'activo' : 'pausado') + '.';

      if (autoBattle && tournament && !tournament.champion) {
        clearPendingBattle();
        runNextBattle();
      }
    });

    fetchPage().catch((error) => {
      status.textContent = error.message;
      grid.innerHTML = '<div class="empty">No fue posible renderizar la Pokédex.</div>';
    });
  </script>
</body>
</html>`

// Frontend serves a small client-side Pokedex page.
func (h *Handler) Frontend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
  _, _ = w.Write([]byte(strings.ReplaceAll(frontendHTML, "__APP_NAME__", h.cfg.App.Name)))
}