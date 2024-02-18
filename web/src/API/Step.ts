type Step = {
    id: string;
    parent: string;
    network: Network;
    results: Results;
}

type RunSpec = {
    steps: number;
    iterations: number;
}