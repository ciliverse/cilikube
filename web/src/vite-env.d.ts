/// <reference types="vite/client" />

declare module '@dagrejs/dagre' {
  const dagre: {
    graphlib: {
      Graph: new () => {
        setDefaultEdgeLabel: (fn: () => object) => void
        setGraph: (attr: Record<string, unknown>) => void
        setNode: (id: string, label: { width: number; height: number }) => void
        setEdge: (source: string, target: string) => void
        node: (id: string) => { x: number; y: number } | undefined
      }
    }
    layout: (g: unknown) => void
  }
  export default dagre
}
