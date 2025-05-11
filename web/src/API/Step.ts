import { Results } from './Results';

type Step = {
    id: string;
    parent: string;
    results: Results;
}

type RunSpec = {
    steps: number;
    iterations: number;
}

export type { Step, RunSpec };