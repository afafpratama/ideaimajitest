const getAuthHeader = () => {
  return {
    'authToken': `${JSON.parse(localStorage.getItem('accountData')).token}`,
    'Access-Control-Allow-Origin': '*'
  }
}

const getOriginHeader = () => {
  return {
    'Access-Control-Allow-Origin': '*'
  }
}

export { getAuthHeader, getOriginHeader }
