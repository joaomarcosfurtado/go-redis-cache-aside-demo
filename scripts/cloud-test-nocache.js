import http from 'k6/http';
import { check } from 'k6';

export const options = {
  vus: 150, // Mesma carga (150 usuários)
  duration: '15s',
  thresholds: {
    http_req_duration: ['p(95)<750'], 
  },
};

export default function () {
  // Gera coordenadas aleatórias para garantir um CACHE MISS
  const randomLat = Math.floor(Math.random() * 1000);
  const randomLon = Math.floor(Math.random() * 1000);

  // Mantemos ?mock=true para simular a lentidão de 500ms (backend lento)
  const url = `https://weather-app-49875691728.us-east1.run.app/weather?lat=${randomLat}&lon=${randomLon}&mock=true`;
  
  const params = { headers: { 'Accept': 'application/json' } };
  
  const res = http.get(url, params);
  
  check(res, {
    'status is 200': (r) => r.status === 200,
    // Verificamos se o header é MISS (para provar que não usou o cache)
    'is MISS': (r) => r.headers['X-Cache'] === 'MISS', 
  });
}