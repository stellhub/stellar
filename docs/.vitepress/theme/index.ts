import DefaultTheme from "vitepress/theme";
import { defineComponent, h, nextTick, onMounted, watch } from "vue";
import { useRoute } from "vitepress";
import "./style.css";

const HOME_ROUTES = new Set(["/", "/zh/"]);

const HomeMatrixBridge = defineComponent({
  name: "HomeMatrixBridge",
  setup() {
    const route = useRoute();

    const syncPageContext = () => {
      const body = document.body;
      if (!body) {
        return;
      }

      delete body.dataset.stellarProductSlug;
      delete body.dataset.stellarProductView;
      delete body.dataset.stellarHome;

      if (HOME_ROUTES.has(route.path)) {
        body.dataset.stellarHome = "true";
      }

      const productMatch = route.path.match(/^\/(?:zh\/)?products\/([^/]+)\/?(.*)$/);
      if (!productMatch) {
        return;
      }

      const [, slug, suffix] = productMatch;
      if (!slug.startsWith("stell")) {
        return;
      }

      body.dataset.stellarProductSlug = slug;
      body.dataset.stellarProductView = suffix ? "subpage" : "home";
    };

    const bindHomeMatrix = () => {
      if (!HOME_ROUTES.has(route.path)) {
        return;
      }

      const rows = Array.from(
        document.querySelectorAll<HTMLTableRowElement>(".stellar-home-matrix tbody tr[data-node]")
      );
      const blueprint = document.querySelector<HTMLObjectElement>("#stellar-home-blueprint");

      if (!rows.length || !blueprint) {
        return;
      }

      const rowMap = new Map(rows.map((row) => [row.dataset.node ?? "", row]));
      let activeNodeId = "";

      const setActiveRow = (nodeId: string) => {
        rows.forEach((row) => {
          row.classList.toggle("is-active", row.dataset.node === nodeId);
        });
      };

      const bindBlueprintDocument = (nodeIdToActivate?: string) => {
        const svgDocument = blueprint.contentDocument;
        if (!svgDocument) {
          return;
        }

        const nodes = Array.from(
          svgDocument.querySelectorAll<SVGGElement>(".bp-node[data-node]")
        );

        const activate = (nodeId: string, source: "table" | "graph" | "init") => {
          if (!nodeId) {
            return;
          }

          activeNodeId = nodeId;
          setActiveRow(nodeId);
          nodes.forEach((node) => {
            node.classList.toggle("active", node.dataset.node === nodeId);
          });

          if (source === "table") {
            blueprint.scrollIntoView({ behavior: "smooth", block: "center" });
          } else if (source === "graph") {
            rowMap.get(nodeId)?.scrollIntoView({ behavior: "smooth", block: "center" });
          }
        };

        nodes.forEach((node) => {
          node.onclick = () => activate(node.dataset.node ?? "", "graph");
        });

        rows.forEach((row) => {
          const nodeId = row.dataset.node ?? "";
          const anchor = row.querySelector<HTMLAnchorElement>("a");

          row.onclick = (event) => {
            const target = event.target as HTMLElement | null;
            if (target?.closest("a")) {
              return;
            }
            activate(nodeId, "table");
          };

          row.onkeydown = (event) => {
            if (event.key !== "Enter" && event.key !== " ") {
              return;
            }
            event.preventDefault();
            activate(nodeId, "table");
          };

          anchor?.addEventListener("click", (event) => {
            if (event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) {
              return;
            }

            event.preventDefault();
            activate(nodeId, "table");
            window.setTimeout(() => {
              window.location.href = anchor.href;
            }, 180);
          });
        });

        if (nodeIdToActivate) {
          activate(nodeIdToActivate, "init");
        }
      };

      if (blueprint.contentDocument) {
        bindBlueprintDocument(activeNodeId || "m01");
      } else {
        blueprint.onload = () => bindBlueprintDocument(activeNodeId || "m01");
      }

      setActiveRow("m01");
    };

    const init = async () => {
      await nextTick();
      syncPageContext();
      window.setTimeout(bindHomeMatrix, 0);
    };

    onMounted(init);
    watch(() => route.path, init);

    return () => null;
  }
});

export default {
  extends: DefaultTheme,
  Layout: () =>
    h(DefaultTheme.Layout, null, {
      "layout-bottom": () => h(HomeMatrixBridge)
    })
};
