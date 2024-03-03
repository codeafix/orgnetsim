type Results = {
    iterations: number;
    colors: Array<Array<number>>;
    conversations: Array<number>;
}

type ResultsCsv = {
    filename: string;
    blob: Blob;
}

export type { Results, ResultsCsv };