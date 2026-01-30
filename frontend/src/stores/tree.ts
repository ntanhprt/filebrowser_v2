import { defineStore } from "pinia";
import { pub as publicApi, files as privateApi } from "@/api";
import { useDebounceFn } from "@vueuse/core";

export interface TreeNode {
    name: string;
    path: string;
    isDir: boolean;
    children: string[];
    loaded: boolean;
    loading: boolean;
}

export const useTreeStore = defineStore("tree", {
    state: () => ({
        nodes: {} as Record<string, TreeNode>,
        expandedPaths: JSON.parse(localStorage.getItem("expandedPaths") || "{}") as Record<string, boolean>,
    }),

    actions: {
        async fetchNode(path: string, isShare: boolean, hash: string | null) {
            const normalizedPath = path === "" ? "/" : path;

            // If already loading, return
            if (this.nodes[path]?.loading) return;

            if (!this.nodes[path]) {
                this.nodes[path] = {
                    name: path.split("/").pop() || "Root",
                    path,
                    isDir: true,
                    children: [],
                    loaded: false,
                    loading: true,
                };
            } else {
                this.nodes[path].loading = true;
            }

            const safePath = normalizedPath.startsWith("/") ? normalizedPath : "/" + normalizedPath;

            try {
                let data;
                if (isShare && hash) {
                    data = await publicApi.fetch(`/share/${hash}${safePath}`);
                } else {
                    data = await privateApi.fetch(`/files${safePath}`);
                }

                const items = data.items.filter((item: any) => item.isDir);
                const childrenPaths: string[] = [];

                items.forEach((item: any) => {
                    let childPath = item.path;
                    if (!childPath.startsWith("/")) childPath = "/" + childPath;
                    childrenPaths.push(childPath);

                    // Initialize child node if not exists
                    if (!this.nodes[childPath]) {
                        this.nodes[childPath] = {
                            name: item.name,
                            path: childPath,
                            isDir: true,
                            children: [],
                            loaded: false,
                            loading: false,
                        };
                    }
                });

                this.nodes[path].children = childrenPaths;
                this.nodes[path].loaded = true;
            } catch (error) {
                console.error("Failed to fetch tree node:", path, error);
            } finally {
                this.nodes[path].loading = false;
            }
        },

        async toggleNode(path: string, isShare: boolean, hash: string | null) {
            this.expandedPaths[path] = !this.expandedPaths[path];
            this.saveToLocalStorage();

            if (this.expandedPaths[path] && !this.nodes[path]?.loaded) {
                await this.fetchNode(path, isShare, hash);
            }
        },

        async ensurePathExpanded(targetPath: string, isShare: boolean, hash: string | null) {
            const parts = targetPath.split("/").filter(p => p);
            let currentPath = "";

            // Expand Root
            if (!this.expandedPaths[""]) {
                this.expandedPaths[""] = true;
                if (!this.nodes[""]?.loaded) await this.fetchNode("", isShare, hash);
            }

            for (const part of parts) {
                currentPath += "/" + part;
                if (currentPath === targetPath) break;

                if (!this.expandedPaths[currentPath]) {
                    this.expandedPaths[currentPath] = true;
                    if (!this.nodes[currentPath]?.loaded) {
                        await this.fetchNode(currentPath, isShare, hash);
                    }
                }
            }
            this.saveToLocalStorage();
        },

        saveToLocalStorage: useDebounceFn(function (this: any) {
            localStorage.setItem("expandedPaths", JSON.stringify(this.expandedPaths));
        }, 1000),

        reset() {
            this.nodes = {};
            // We might want to keep expandedPaths for persistence
        }
    },
});
