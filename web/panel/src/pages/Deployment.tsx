import { Anchor, Center, Group, Paper, Title, Text, Stack, HoverCard } from "@mantine/core";
import { DeploymentDetailResp } from "../utils/types";
import { deploymentsMeta } from "../data/deployment";
import { useEffect, useState } from "react";
import { Github } from "lucide-react";
import Loading from "./Loading";

/*
  ToDo:
  Commit hash should be a link to the commit on GitHub
  Show other meta data in hover card
*/

export default function Deployment() {
  const [deploymentData, setDeploymentData] = useState<DeploymentDetailResp>();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // ToDo: Fetch actual deployment data from API
    // ToDo: Apply some data transformation

    setLoading(false);
    setDeploymentData({
      meta: {
        StartAt: "2021-09-01 12:00:00",
        EndAt: "2021-09-01 12:00:00",
        CommitHash: "0571234567890",
        CommitMsg: "Initial Commit",
        Repo: "Username/Repo",
        Status: "Success",
        Branch: "main",
        ID: 90,
      },
      message: "Deployment details fetched successfully",
      ok: true,
      logs: []
    });
  }, []);

  if (loading) {
    return <Loading />;
  }

  const { meta } = deploymentData!;
  
  return (
    <>
      <Center>
        <Paper w="90%" shadow="sm" my={20} p={20}>
          <Anchor href="http://github.com" target="_blank">
            <Group gap={8} pb={10}>
              <Github />
              <Title order={2}>
                {meta.Repo}
              </Title>
            </Group>
          </Anchor>
          <Stack gap={5}>
            {
              deploymentsMeta.map(({ key, hovText, icon }) => {
                return (
                  <HoverCard key={key} closeDelay={0} position="bottom-start" shadow="md" withArrow>
                    <HoverCard.Target>
                      <Group w="fit-content" gap={8}>
                        {icon}
                        <Text>{meta[key]}</Text>
                      </Group>
                    </HoverCard.Target>
                    <HoverCard.Dropdown>
                      <Text>{hovText}</Text>
                    </HoverCard.Dropdown>
                  </HoverCard>
                )
              })
            }
          </Stack>
        </Paper>
      </Center>

      <Center>
        <Paper w="90%" shadow="sm" my={20} p={20}>
          <Title order={2}>
            Logs
          </Title>
          <Text>
            No logs available
          </Text>
          {/* ToDo: Add timeline element for logs */}
        </Paper>
      </Center>
    </>
  )
}
