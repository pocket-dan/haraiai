export const siteName = "haraiai (払い合い)";
export const siteOrigin = "https://haraiai.com";

export const pageUrlOf = (path: string) => {
  return path === "/" ? siteOrigin : siteOrigin + path;
};
