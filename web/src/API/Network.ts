type Network = {
    links: Array<Link>
    nodes: Array<AgentState>;
    maxColors: number;
}

type Link = {
    source: string;
    target: string;
    strength: number;
    length: number;
}

type AgentState = {
    id: string;
    color: number;
    susceptability: number;
    influence: number;
    contrariness: number;
    change: number;
    type: string;
    fx: number;
    fy: number;
}
