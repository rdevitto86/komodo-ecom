export async function load({ setHeaders }) {
  setHeaders({
    'cache-control': 'public, max-age=300, s-maxage=600'
  });
  
  return {
    seed: Math.random()
  };
}
