import { ActionIcon, Card, Flex, Group, Image, Input, Title, useMantineColorScheme } from '@mantine/core';
import { Search, SunMoon } from "lucide-react";
import { useMeta } from '../state/use-meta';
import { useAuth } from '../state/use-auth';

export default function Header() {
  const { toggleColorScheme } = useMantineColorScheme();
  const { pageTitle } = useMeta();
  const { loggedIn } = useAuth();

  return (
    <Card
      withBorder
      radius={0}
      shadow='md'
      style={{ position: 'sticky', top: 0, zIndex: 1000 }}
    >
      <Flex
        p={0}
        justify="space-between"
      >
        <Group gap="xs">
          <Image
            w={{ xs: 36, sm: 54 }}
            src="/media/d2m.png"
            alt="Logo of D2M"
          />
          <Title
            order={2}
            fz={{ xs: 18, sm: 24 }}
          >
            {pageTitle}
          </Title>
        </Group>

        <Group gap="xs">
          {loggedIn && (
            <Input
              readOnly
              radius={50}
              w={150}
              variant='filled'
              placeholder='Search...'
              leftSection={<Search />}
            />
          )}
          <ActionIcon
            aria-label='Theme switch icon'
            onClick={toggleColorScheme}
            variant='subtle'
            size="xl"
          >
            <SunMoon size={28} />
          </ActionIcon>
        </Group>
      </Flex>
    </Card>
  )
}
