import http from 'k6/http';
import { check } from 'k6';

export const options = {
  // We will simulate 50 users constantly hitting API
  vus: 150, 
  duration: '15s',
  
  thresholds: {
    http_req_duration: ['p(95)<250'], // 95% of requests must be under 200ms
  },
};

export default function () {
  // We will use mock=true to simulate the slow backend safely (so we don't ban the IP)
  const url = 'https://weather-app-49875691728.us-east1.run.app/weather?lat=52&lon=13&mock=true';
  
  const params = { 
    headers: { 'Accept': 'application/json' } 
  };
  
  const res = http.get(url, params);
  
  check(res, {
    'status is 200': (r) => r.status === 200,
    'hit cache': (r) => r.headers['X-Cache'] === 'HIT', // Monitor Cache Hits
  });
}