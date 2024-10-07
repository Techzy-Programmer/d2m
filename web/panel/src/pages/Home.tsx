import useFetch from "../hooks/useFetch";
import HomeMeta from "./chunks/HomeMeta";
import { useEffect, useState } from "react";
import { useMeta } from "../state/use-meta";
import { DeploymentsResp } from "../utils/types";
import { Center, Paper, Table } from "@mantine/core";
import { isFetchSuccess, showToast } from "../utils/general";

function formatDate(date: string) {
  return (new Date(date).toLocaleString('en-GB', {
    weekday: 'short',
    year: '2-digit',
    month: 'short',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    timeZone: "IST",
    hour12: true
  })) + " IST";
}

export default function Home() {
  const [deployments, setDeployments] = useState<DeploymentsResp["deployments"]>([]);
  const { metadata } = useMeta();
  const fetchData = useFetch();

  useEffect(() => {
    (async () => {
      const resp = await fetchData<DeploymentsResp>("/api/mg/get-deployments");
      
      if (!isFetchSuccess(resp)) {
        return showToast({
          message: resp.error || "Something went wrong",
          title: "Failed To Fetch",
          status: "issue",
        });
      }

      const { code, data } = resp;

      if (code !== 200) {
        return showToast({
          title: "Server's mood is off",
          message: data.message,
          status: "issue",
        });
      }

      setDeployments(data.deployments);
    })();
  }, []);

  const rows = deployments.map((dep) => (
    <Table.Tr key={dep.ID}>
      <Table.Td>{dep.ID}</Table.Td>
      <Table.Td>{dep.Repo}</Table.Td>
      <Table.Td>{dep.Branch}</Table.Td>
      <Table.Td>{dep.CommitMsg}</Table.Td>
      <Table.Td>{formatDate(dep.StartAt)}</Table.Td>
      <Table.Td>{formatDate(dep.EndAt)}</Table.Td>
      <Table.Td>{dep.Status}</Table.Td>
    </Table.Tr>
  ));

  return (
    <>
      <HomeMeta {...metadata} />
      <Center>
        <Paper
          mt={10}
          w="90%"
          p="md"
        >
          <Table.ScrollContainer minWidth={1200}>
            <Table
              striped
              withTableBorder
              highlightOnHover
              withColumnBorders
            >
              <Table.Thead>
                <Table.Tr>
                  <Table.Th>ID</Table.Th>
                  <Table.Th>Repo</Table.Th>
                  <Table.Th>Branch</Table.Th>
                  <Table.Th>Commit Message</Table.Th>
                  <Table.Th>Start At</Table.Th>
                  <Table.Th>End At</Table.Th>
                  <Table.Th>Status</Table.Th>
                </Table.Tr>
              </Table.Thead>
              <Table.Tbody>{rows}</Table.Tbody>
            </Table>
          </Table.ScrollContainer>
        </Paper>
      </Center>
    </>
  )
}
