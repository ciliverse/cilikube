import axios from 'axios';

const baseURL = 'http://127.0.0.1:8080';

async function testAPI() {
  try {
    console.log('Testing backend API connection...');
    
    // Test cluster list endpoint
    console.log('1. Testing /api/v1/clusters...');
    const clustersResponse = await axios.get(`${baseURL}/api/v1/clusters`);
    console.log('Clusters response:', clustersResponse.data);
    
    // Test active cluster endpoint
    console.log('2. Testing /api/v1/clusters/active...');
    try {
      const activeResponse = await axios.get(`${baseURL}/api/v1/clusters/active`);
      console.log('Active cluster response:', activeResponse.data);
    } catch (error) {
      console.log('Active cluster error:', error.response?.data || error.message);
    }
    
    // Test summary endpoint
    console.log('3. Testing /api/v1/summary/resources...');
    try {
      const summaryResponse = await axios.get(`${baseURL}/api/v1/summary/resources`);
      console.log('Summary response:', summaryResponse.data);
    } catch (error) {
      console.log('Summary error:', error.response?.data || error.message);
    }
    
  } catch (error) {
    console.error('API test failed:', error.message);
    if (error.code === 'ECONNREFUSED') {
      console.error('Backend server is not running on port 8080');
    }
  }
}

testAPI();