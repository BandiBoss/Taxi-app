const axiosMock = jest.createMockFromModule('axios');

axiosMock.create = jest.fn(() => axiosMock);

module.exports = axiosMock;