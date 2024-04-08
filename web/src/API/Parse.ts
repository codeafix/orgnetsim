type ParseOptions = {
    identifier: number;
    parent: number;
    name: number;
    regex: Regex;
    delimiter: string;
    payload: string;
}

type Regex = {
    [key: string]: string;
}

export type { ParseOptions, Regex };