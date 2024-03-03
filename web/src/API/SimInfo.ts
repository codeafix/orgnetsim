import { NetworkOptions } from './NetworkOptions';

type SimInfo = {
    id: string;
    name: string;
    description: string;
    steps: Array<string>;
    options: NetworkOptions;
}

export type { SimInfo };