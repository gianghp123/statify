export type CamelCase<S extends string> = S extends `${infer P1}_${infer P2}${infer P3}`
  ? `${Lowercase<P1>}${Uppercase<P2>}${CamelCase<P3>}`
  : Lowercase<S>;

export type Camelize<T> = {
  [K in keyof T as CamelCase<string & K>]: T[K] extends object
    ? Camelize<T[K]>
    : T[K];
};

function toCamel(s: string) {
  return s.replace(/([-_][a-z])/gi, ($1) => {
    return $1.toUpperCase().replace("-", "").replace("_", "");
  });
}

export function snakeToCamel<T>(obj: any): T {
  if (Array.isArray(obj)) {
    return obj.map((v) => snakeToCamel(v)) as any;
  } else if (obj !== null && obj.constructor === Object) {
    return Object.keys(obj).reduce(
      (result, key) => ({
        ...result,
        [toCamel(key)]: snakeToCamel(obj[key]),
      }),
      {}
    ) as any;
  }
  return obj;
}
