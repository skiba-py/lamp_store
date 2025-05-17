import { Box, Flex, Text } from '@chakra-ui/react';
import { FaRegLightbulb } from 'react-icons/fa';
import { Link } from 'react-router-dom';

export default function Footer() {
  return (
    <Box as="footer" bg="white" py={3}>
      <Flex direction="column" align="center" maxW="1200px" mx="auto">
        <Flex w="100%" justify="center" mb={2}>
          <Link to="/" style={{ display: 'flex', alignItems: 'center', textDecoration: 'none' }}>
            <Text fontWeight="bold">LAMPS</Text>
            <FaRegLightbulb />
          </Link>
        </Flex>
        <Flex w="100%" justify="center" fontSize="sm" color="gray.500">
          <Text>© 2025 </Text>
        </Flex>
      </Flex>
    </Box>
  );
}
