---
layout: home

hero:
  name: Stell Hub
  text: Platform Infrastructure Documentation
  tagline: English-first reference material for the Stell Hub middleware stack, product architecture, and engineering design notes.
  actions:
    - theme: brand
      text: Explore Products
      link: /products/
    - theme: alt
      text: Explore Topics
      link: /topics/
---

<div class="stellar-home-wide">
  <section class="stellar-home-matrix">
    <h2>Middleware Capability Matrix</h2>
    <div class="stellar-home-capability-legend" aria-label="Capability legend">
      <span class="stellar-home-capability-legend__item is-access">
        <span class="stellar-home-capability-legend__swatch"></span>
        <span>ACCESS / SECURITY</span>
      </span>
      <span class="stellar-home-capability-legend__item is-runtime">
        <span class="stellar-home-capability-legend__swatch"></span>
        <span>GOVERNANCE / RUNTIME</span>
      </span>
      <span class="stellar-home-capability-legend__item is-foundation">
        <span class="stellar-home-capability-legend__swatch"></span>
        <span>FOUNDATION</span>
      </span>
      <span class="stellar-home-capability-legend__item is-observe">
        <span class="stellar-home-capability-legend__swatch"></span>
        <span>OBSERVABILITY / ALERT</span>
      </span>
    </div>
    <div class="stellar-home-matrix__table">
      <table>
        <colgroup>
          <col class="stellar-home-matrix__col-id">
          <col class="stellar-home-matrix__col-domain">
          <col class="stellar-home-matrix__col-origin">
          <col class="stellar-home-matrix__col-name">
          <col class="stellar-home-matrix__col-logic">
        </colgroup>
        <thead>
          <tr>
            <th><span class="stellar-home-matrix__tag">ID</span></th>
            <th><span class="stellar-home-matrix__tag">DOMAIN</span></th>
            <th><span class="stellar-home-matrix__tag">ORIGIN</span></th>
            <th><span class="stellar-home-matrix__tag">STELL NAME</span></th>
            <th><span class="stellar-home-matrix__tag">RATIONALE</span></th>
          </tr>
        </thead>
        <tbody>
          <tr data-node="m01" data-group="foundation" tabindex="0">
            <td><span class="stellar-home-matrix__id">M01</span></td>
            <td><strong>REGISTRY CENTER</strong></td>
            <td>StarMap</td>
            <td><strong><a href="/products/stellmap/">Stellmap</a></strong></td>
            <td>Stell + Map, a service coordinate map for discovery and routing.</td>
          </tr>
          <tr data-node="m02" data-group="foundation" tabindex="0">
            <td><span class="stellar-home-matrix__id">M02</span></td>
            <td><strong>CONFIGURATION CENTER</strong></td>
            <td>Nebula</td>
            <td><strong><a href="/products/stellnula/">Stellnula</a></strong></td>
            <td>A condensed form of Nebula, emphasizing widely distributed configuration presence.</td>
          </tr>
          <tr data-node="m03" data-group="observe" tabindex="0">
            <td><span class="stellar-home-matrix__id">M03</span></td>
            <td><strong>DISTRIBUTED TRACING</strong></td>
            <td>StarTrace</td>
            <td><strong><a href="/products/stelltrace/">Stelltrace</a></strong></td>
            <td>Stell + Trace, tracking request paths across a distributed landscape.</td>
          </tr>
          <tr data-node="m04" data-group="runtime" tabindex="0">
            <td><span class="stellar-home-matrix__id">M04</span></td>
            <td><strong>SERVICE GOVERNANCE</strong></td>
            <td>Orbit</td>
            <td><strong><a href="/products/stellorbit/">Stellorbit</a></strong></td>
            <td>Services stay on predictable operational trajectories, like bodies on an orbit.</td>
          </tr>
          <tr data-node="m05" data-group="runtime" tabindex="0">
            <td><span class="stellar-home-matrix__id">M05</span></td>
            <td><strong>RATE LIMITING AND CIRCUIT BREAKING</strong></td>
            <td>Pulsar</td>
            <td><strong><a href="/products/stellpulse/">Stellpulse</a></strong></td>
            <td>Stell + Pulse, focused on sensing and controlling live traffic rhythm.</td>
          </tr>
          <tr data-node="m06" data-group="runtime" tabindex="0">
            <td><span class="stellar-home-matrix__id">M06</span></td>
            <td><strong>JOB SCHEDULING</strong></td>
            <td>Astrolabe</td>
            <td><strong><a href="/products/stellabe/">Stellabe</a></strong></td>
            <td>A shortened form of Astrolabe, referencing coordinated timing and positioning.</td>
          </tr>
          <tr data-node="m07" data-group="foundation" tabindex="0">
            <td><span class="stellar-home-matrix__id">M07</span></td>
            <td><strong>DISTRIBUTED LOCK</strong></td>
            <td>Singularity</td>
            <td><strong><a href="/products/stellpoint/">Stellpoint</a></strong></td>
            <td>A single decisive point, fitting exclusive ownership and coordination semantics.</td>
          </tr>
          <tr data-node="m08" data-group="access" tabindex="0">
            <td><span class="stellar-home-matrix__id">M08</span></td>
            <td><strong>GATEWAY</strong></td>
            <td>Event Horizon</td>
            <td><strong><a href="/products/stellgate/">Stellgate</a></strong></td>
            <td>A gateway is the entry boundary, like a horizon separating external and internal traffic.</td>
          </tr>
          <tr data-node="m09" data-group="runtime" tabindex="0">
            <td><span class="stellar-home-matrix__id">M09</span></td>
            <td><strong>MESSAGE QUEUE</strong></td>
            <td>CometFlow</td>
            <td><strong><a href="/products/stellflow/">Stellflow</a></strong></td>
            <td>Ordered data streams that travel like a comet trail.</td>
          </tr>
          <tr data-node="m10" data-group="observe" tabindex="0">
            <td><span class="stellar-home-matrix__id">M10</span></td>
            <td><strong>LOG PLATFORM</strong></td>
            <td>Spectrum</td>
            <td><strong><a href="/products/stellspec/">Stellspec</a></strong></td>
            <td>System behavior is inspected through its signal spectrum.</td>
          </tr>
          <tr data-node="m11" data-group="observe" tabindex="0">
            <td><span class="stellar-home-matrix__id">M11</span></td>
            <td><strong>METRICS PLATFORM</strong></td>
            <td>Constellation</td>
            <td><strong><a href="/products/stellcon/">Stellcon</a></strong></td>
            <td>Metrics points connect into recognizable constellations for operational reasoning.</td>
          </tr>
          <tr data-node="m12" data-group="observe" tabindex="0">
            <td><span class="stellar-home-matrix__id">M12</span></td>
            <td><strong>ALERTING PLATFORM</strong></td>
            <td>NovaSignal</td>
            <td><strong><a href="/products/stellvox/">Stellvox</a></strong></td>
            <td>Vox means voice or message, matching an alerting system's notification role.</td>
          </tr>
          <tr data-node="m13" data-group="access" tabindex="0">
            <td><span class="stellar-home-matrix__id">M13</span></td>
            <td><strong>ZERO TRUST PLATFORM</strong></td>
            <td>StarShield</td>
            <td><strong><a href="/products/stellguard/">Stellguard</a></strong></td>
            <td>Dynamic guard semantics suit continuous verification better than a static shield metaphor alone.</td>
          </tr>
          <tr data-node="m14" data-group="foundation" tabindex="0">
            <td><span class="stellar-home-matrix__id">M14</span></td>
            <td><strong>KEY MANAGEMENT CENTER</strong></td>
            <td>StarKey</td>
            <td><strong><a href="/products/stellkey/">Stellkey</a></strong></td>
            <td>A direct and forceful name for secret and certificate management responsibilities.</td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</div>

<div class="stellar-home-wide">
  <section class="stellar-home-graph">
    <h2>Dependency Blueprint</h2>
    <p>Use this view to read the stack by platform layers, core contracts, and observability feedback paths.</p>
    <object
      id="stellar-home-blueprint"
      class="stellar-home-blueprint"
      data="/core-middleware-matrix-dependency.svg"
      type="image/svg+xml"
      aria-label="Middleware dependency blueprint"
    ></object>
  </section>
</div>

## Contact

- GitHub: [https://github.com/stellhub](https://github.com/stellhub)
- Email: [xiaoyaoyunlian@gmail.com](mailto:xiaoyaoyunlian@gmail.com)
- Backup email: [chenwenlong_java@163.com](mailto:chenwenlong_java@163.com)
