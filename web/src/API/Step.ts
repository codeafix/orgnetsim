import { Network } from './Network';
import { Results } from './Results';

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

export type { Step, RunSpec };